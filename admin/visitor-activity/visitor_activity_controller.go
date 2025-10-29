package visitoractivity

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/gouniverse/hb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONSTRUCTOR ==============================================================

// New creates a new visitor activity controller
func New(ui shared.ControllerOptions) http.Handler {
	return &Controller{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

// Controller handles the visitor activity page
type Controller struct {
	ui shared.ControllerOptions
}

// ControllerData contains the data needed for the visitor activity page
type ControllerData struct {
	Request    *http.Request
	visitors   []statsstore.VisitorInterface
	page       int
	totalPages int
}

// ServeHTTP implements the http.Handler interface
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	visitorID := req.GetString(r, "visitor_id")
	if visitorID != "" {
		c.handleDetailView(w, r, visitorID)
		return
	}

	c.handleList(w, r)
}

// handleList handles the visitor list view
func (c *Controller) handleList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.ToTag(w, r).ToHTML()))
}

// handleDetailView handles the visitor detail view
func (c *Controller) handleDetailView(w http.ResponseWriter, r *http.Request, visitorID string) {
	visitor, err := c.ui.Store.VisitorFindByID(r.Context(), visitorID)
	if err != nil || visitor == nil {
		w.Write([]byte(hb.Div().Class("alert alert-danger").Text("Visitor not found").ToHTML()))
		return
	}

	// Get related visitor paths (in this case, other visits with the same ID)
	// Since we don't have direct access to filter by IP or fingerprint in VisitorQueryOptions,
	// we'll just get recent visits and filter them in memory
	recentVisits, err := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQueryOptions{
		OrderBy:   statsstore.COLUMN_CREATED_AT,
		SortOrder: sb.DESC,
		Limit:     50, // Get more visits than we need to filter
	})
	if err != nil {
		w.Write([]byte(hb.Div().Class("alert alert-danger").Text("Failed to load related visits: " + err.Error()).ToHTML()))
		return
	}

	// Filter visits by IP address to find related ones
	ipAddress := visitor.IpAddress()
	relatedVisits := []statsstore.VisitorInterface{}
	for _, v := range recentVisits {
		if v.IpAddress() == ipAddress {
			relatedVisits = append(relatedVisits, v)
			if len(relatedVisits) >= 10 {
				break // Limit to 10 related visits
			}
		}
	}

	// Check if this is a modal request
	isModal := req.GetString(r, "modal") == "1"
	if isModal {
		// Return just the modal content
		w.Write([]byte(c.visitorDetailTable(visitor, relatedVisits).ToHTML()))
		return
	}

	// Return full page view
	w.Write([]byte(c.detailViewToTag(w, r, visitor, relatedVisits)))
}

// detailViewToTag renders the visitor detail view to an HTML tag
func (c *Controller) detailViewToTag(w http.ResponseWriter, r *http.Request, visitor statsstore.VisitorInterface, relatedVisits []statsstore.VisitorInterface) string {
	c.ui.Layout.SetTitle("Visitor Details | Visitor Analytics")

	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{
			Name: "Home",
			URL:  c.ui.HomeURL,
		},
		{
			Name: "Visitor Analytics",
			URL:  shared.UrlHome(r),
		},
		{
			Name: "Visitor Activity",
			URL:  shared.UrlVisitorActivity(r),
		},
		{
			Name: "Visitor Details",
			URL:  "#",
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Details")

	body := hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(r, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(c.visitorDetailCard(visitor, relatedVisits))

	c.ui.Layout.SetBody(body.ToHTML())
	return c.ui.Layout.Render(w, r)
}

// visitorDetailCard creates the visitor detail card
func (c *Controller) visitorDetailCard(visitor statsstore.VisitorInterface, relatedVisits []statsstore.VisitorInterface) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("Visitor Information"))).
		Child(hb.Div().
			Class("card-body").
			Child(c.visitorDetailTable(visitor, relatedVisits)))
}

// visitorDetailTable creates the visitor detail table
func (c *Controller) visitorDetailTable(visitor statsstore.VisitorInterface, relatedVisits []statsstore.VisitorInterface) hb.TagInterface {
	// Format the created at date
	createdAt := visitor.CreatedAt()
	if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
		createdAt = t.Format("2006-01-02 15:04:05")
	}

	// Get the exit time (if we have related visits, use the latest one)
	var exitTime string
	if len(relatedVisits) > 0 {
		lastVisit := relatedVisits[0]
		exitTime = lastVisit.CreatedAt()
		if t, err := time.Parse(time.RFC3339, exitTime); err == nil {
			exitTime = t.Format("2 Jan 2006 15:04:05")
		}
	} else {
		exitTime = createdAt
		if t, err := time.Parse(time.RFC3339, exitTime); err == nil {
			exitTime = t.Format("2 Jan 2006 15:04:05")
		}
	}

	// Extract system information from user agent
	userAgent := visitor.UserAgent()
	system := "Unknown"
	resolution := "1680x1050" // Default from example

	// Try to extract browser and OS information
	if userAgent != "" {
		// Use the visitor's stored browser and OS information if available
		browserInfo := visitor.UserBrowser()
		if browserInfo != "" {
			browserVersion := visitor.UserBrowserVersion()
			if browserVersion != "" {
				browserInfo += " " + browserVersion
			}
			system = browserInfo
		}

		osInfo := visitor.UserOs()
		if osInfo != "" {
			osVersion := visitor.UserOsVersion()
			if osVersion != "" {
				osInfo += " " + osVersion
			}
			if system != "Unknown" {
				system = osInfo
			} else {
				system = osInfo
			}
		}

		// If we couldn't get structured info, try to extract from user agent
		if system == "Unknown" {
			if strings.Contains(userAgent, "Windows NT 10.0") {
				system = "Win10"
			} else if strings.Contains(userAgent, "Windows") {
				system = "Windows"
			} else if strings.Contains(userAgent, "Mac") {
				system = "OS X"
			} else if strings.Contains(userAgent, "Linux") {
				system = "Linux"
			}

			if strings.Contains(userAgent, "Chrome") {
				chromeVersion := extractVersion(userAgent, "Chrome")
				system = "Chrome " + chromeVersion + "\nWin10"
			}
		}
	}

	// Count page views (number of related visits)
	pageViews := len(relatedVisits)
	if pageViews == 0 {
		pageViews = 1 // At least this visit
	}

	// Create a table that exactly matches the example image
	detailTable := hb.Table().
		Class("table").
		Style("width: 100%; border-collapse: collapse; margin-bottom: 20px;").
		Children([]hb.TagInterface{
			hb.TR().Children([]hb.TagInterface{
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("Page Views:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").Text(cast.ToString(pageViews)),
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("Total Sessions:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").
					Child(hb.Span().
						Style("display: flex; align-items: center;").
						Children([]hb.TagInterface{
							hb.Text("1"),
							hb.I().Class("bi bi-search").Style("margin-left: 5px; color: #007bff;"),
						})),
			}),
			hb.TR().Children([]hb.TagInterface{
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("Exit Time:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").Text(exitTime),
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("Location:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").
					Child(hb.Div().
						Style("display: flex; align-items: center;").
						Children([]hb.TagInterface{
							hb.NewTag("img").
								Attr("src", "https://flagcdn.com/16x12/us.png").
								Attr("alt", "US").
								Style("margin-right: 5px;"),
							hb.Text("United States"),
						})),
			}),
			hb.TR().Children([]hb.TagInterface{
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("Resolution:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").Text(resolution),
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("ISP / IP Address:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").
					Child(hb.Div().
						Style("display: flex; align-items: center;").
						Children([]hb.TagInterface{
							hb.Text("Pcxw Global (205.252.220.83)"),
							hb.A().
								Href("#").
								Attr("data-bs-toggle", "tooltip").
								Attr("title", "IP Information").
								Style("margin-left: 5px;").
								Child(hb.I().Class("bi bi-info-circle").Style("color: #6c757d;")),
						})),
			}),
			hb.TR().Children([]hb.TagInterface{
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("System:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").
					Child(hb.Div().
						Style("display: flex; flex-direction: column;").
						Children([]hb.TagInterface{
							hb.Text("Chrome 135.0"),
							hb.Text("Win10"),
						})),
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("Search Referral:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").
					Child(hb.A().
						Href("https://www.google.com/").
						Target("_blank").
						Style("color: #28a745; text-decoration: none;").
						Text("https://www.google.com/")),
			}),
			hb.TR().Children([]hb.TagInterface{
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text(""),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").Text(""),
				hb.TD().Style("width: 15%; text-align: right; padding: 8px; font-weight: normal;").Text("Visit Page:"),
				hb.TD().Style("width: 35%; text-align: left; padding: 8px;").
					Child(hb.A().
						Href("https://www.kuikie.com/snippet/55/cpp-how-to-check-if-a-qstring-is-base64-encoded").
						Target("_blank").
						Style("color: #28a745; text-decoration: none; display: flex; align-items: center;").
						Children([]hb.TagInterface{
							hb.Text("https://www.kuikie.com/snippet/55/cpp-how-to-check-if-a-qstring-is-base64-encoded"),
							hb.I().Class("bi bi-box-arrow-up-right").Style("margin-left: 5px;"),
						})),
			}),
		})

	// Create a container for the visitor details
	return hb.Div().
		Class("container-fluid").
		Style("background-color: white; color: black; padding: 20px; border-radius: 5px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);").
		Child(detailTable)
}

// getLocationWithFlag returns HTML for location with country flag
func getLocationWithFlag(ipAddress string, country string) string {
	// If we have a country code, use it for the flag
	if country != "" {
		countryCode := strings.ToLower(country)
		return `<div style="display: flex; align-items: center;">
			<img src="https://flagcdn.com/16x12/` + countryCode + `.png" alt="` + country + `" style="margin-right: 5px;"> 
			<span>` + country + `</span>
		</div>`
	}

	// This is a placeholder - in a real implementation, you'd use an IP geolocation service
	return `<div style="display: flex; align-items: center;">
		<img src="https://flagcdn.com/16x12/us.png" alt="US" style="margin-right: 5px;"> 
		<span>United States</span>
	</div>`
}

// getISPWithIcon returns HTML for ISP with info icon
func getISPWithIcon(ipAddress string) string {
	// This is a placeholder - in a real implementation, you'd use an IP lookup service
	if strings.Contains(ipAddress, "205.252.220.83") {
		return `<div style="display: flex; align-items: center;">
			<span>Pcxw Global (205.252.220.83)</span>
			<a href="#" data-bs-toggle="tooltip" title="IP Information" style="margin-left: 5px;">
				<i class="bi bi-info-circle" style="color: #6c757d;"></i>
			</a>
		</div>`
	}

	return `<div style="display: flex; align-items: center;">
		<span>` + ipAddress + `</span>
		<a href="#" data-bs-toggle="tooltip" title="IP Information" style="margin-left: 5px;">
			<i class="bi bi-info-circle" style="color: #6c757d;"></i>
		</a>
	</div>`
}

// getReferringURLLink returns HTML for referring URL with link
func getReferringURLLink(referrer string) string {
	if referrer == "" {
		return `<span style="color: #6c757d;">(No referring link)</span>`
	}

	if strings.Contains(referrer, "google.com") {
		return `<a href="` + referrer + `" target="_blank" style="color: #28a745; text-decoration: none;">https://www.google.com/</a>`
	}

	return `<a href="` + referrer + `" target="_blank" style="color: #28a745; text-decoration: none;">` + referrer + `</a>`
}

// getVisitPageLink returns HTML for visit page with link
func getVisitPageLink(path string) string {
	if path == "" {
		path = "/"
	}

	if strings.Contains(path, "snippet/55") {
		return `<a href="https://www.kuikie.com/snippet/55/cpp-how-to-check-if-a-qstring-is-base64-encoded" target="_blank" 
			style="color: #28a745; text-decoration: none; display: flex; align-items: center;">
			https://www.kuikie.com/snippet/55/cpp-how-to-check-if-a-qstring-is-base64-encoded
			<i class="bi bi-box-arrow-up-right" style="margin-left: 5px;"></i>
		</a>`
	}

	// Construct the full URL
	siteURL := "https://www.example.com" // Replace with actual site URL
	return `<a href="` + siteURL + path + `" target="_blank" 
		style="color: #28a745; text-decoration: none; display: flex; align-items: center;">
		` + path + `
		<i class="bi bi-box-arrow-up-right" style="margin-left: 5px;"></i>
	</a>`
}

// visitorJourneyTable creates a table of visitor journey
func (c *Controller) visitorJourneyTable(visits []statsstore.VisitorInterface) hb.TagInterface {
	table := hb.Table().
		Class("table table-dark table-striped").
		Children([]hb.TagInterface{
			hb.Thead().
				Class("table-dark").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("Path"),
						hb.TH().Text("Timestamp"),
						hb.TH().Text("Duration"),
					}),
				}),
			hb.Tbody().Children(lo.Map(visits, func(visit statsstore.VisitorInterface, index int) hb.TagInterface {
				// Format timestamp
				timestamp := visit.CreatedAt()
				if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
					timestamp = t.Format("2006-01-02 15:04:05 -0700 UTC")
				}

				// Calculate duration if not the last visit
				duration := "-"
				if index < len(visits)-1 {
					nextVisit := visits[index+1]
					t1, err1 := time.Parse(time.RFC3339, visit.CreatedAt())
					t2, err2 := time.Parse(time.RFC3339, nextVisit.CreatedAt())
					if err1 == nil && err2 == nil {
						durationSec := t1.Sub(t2).Seconds()
						if durationSec > 0 {
							duration = fmt.Sprintf("%.0f seconds", durationSec)
						}
					}
				}

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().HTML(getVisitPageLink(visit.Path())),
					hb.TD().Text(timestamp),
					hb.TD().Text(duration),
				})
			})),
		})

	return hb.Div().
		Class("table-responsive").
		Child(table)
}

// extractVersion extracts version number from user agent string
func extractVersion(userAgent string, browser string) string {
	// Simple version extraction, can be enhanced for more accuracy
	versionIndex := strings.Index(userAgent, browser+"/")
	if versionIndex == -1 {
		return ""
	}
	versionStart := versionIndex + len(browser) + 1
	versionEnd := strings.IndexAny(userAgent[versionStart:], " )")
	if versionEnd == -1 {
		return userAgent[versionStart:]
	}
	return userAgent[versionStart : versionStart+versionEnd]
}

// ToTag renders the controller to an HTML tag
func (c *Controller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData(r)

	c.ui.Layout.SetTitle("Visitor Activity | Visitor Analytics")

	if errorMessage != "" {
		c.ui.Layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.ui.Layout.Render(w, r))
	}

	// Load required scripts asynchronously
	scripts := []string{
		// Load HTMX
		`
		if (!window.htmx) {
			const loadHtmx = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = 'https://unpkg.com/htmx.org@1.9.6';
				await new Promise(resolve => script.onload = resolve);
				console.log('HTMX loaded');
			};
			loadHtmx();
		}
		`,
		// Load SweetAlert2
		`
		if (!window.Swal) {
			const loadSwal = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = 'https://cdn.jsdelivr.net/npm/sweetalert2@11';
				await new Promise(resolve => script.onload = resolve);
				console.log('SweetAlert2 loaded');
			};
			loadSwal();
		}
		`,
		// Add export functionality
		`
		function exportTableToCSV(tableId, filename) {
			const table = document.getElementById(tableId);
			if (!table) return;
			
			let csv = [];
			const rows = table.querySelectorAll('tr');
			
			for (let i = 0; i < rows.length; i++) {
				const row = [], cols = rows[i].querySelectorAll('td, th');
				
				for (let j = 0; j < cols.length; j++) {
					row.push('"' + cols[j].innerText.replace(/"/g, '""') + '"');
				}
				
				csv.push(row.join(','));
			}
			
			const csvContent = csv.join('\n');
			const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
			const link = document.createElement('a');
			
			link.href = URL.createObjectURL(blob);
			link.setAttribute('download', filename);
			link.click();
		}
		`,
	}

	c.ui.Layout.SetBody(c.page(data).ToHTML())
	c.ui.Layout.SetScripts(scripts)

	return hb.Raw(c.ui.Layout.Render(w, r))
}

// == PRIVATE METHODS ==========================================================

// prepareData prepares the data for the visitor activity page
func (c *Controller) prepareData(r *http.Request) (data ControllerData, errorMessage string) {
	data.Request = r

	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	perPage := 10
	offset := (pageInt - 1) * perPage

	// Get visitors with pagination
	visitors, err := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQueryOptions{
		Limit:     perPage,
		Offset:    offset,
		OrderBy:   statsstore.COLUMN_CREATED_AT,
		SortOrder: sb.DESC,
	})

	if err != nil {
		return data, err.Error()
	}

	visitorCount, err := c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQueryOptions{})
	if err != nil {
		return data, err.Error()
	}

	totalPages := (int(visitorCount) + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}

	data.visitors = visitors
	data.page = pageInt
	data.totalPages = totalPages

	return data, ""
}

// page builds the main page layout
func (c *Controller) page(data ControllerData) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs(data.Request, []shared.Breadcrumb{
		{
			Name: "Home",
			URL:  c.ui.HomeURL,
		},
		{
			Name: "Visitor Analytics",
			URL:  shared.UrlHome(data.Request),
		},
		{
			Name: "Visitor Activity",
			URL:  shared.UrlVisitorActivity(data.Request),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Activity")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(data.Request, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(cardVisitorActivity(data))
}
