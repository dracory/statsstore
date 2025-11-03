package shared

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/hb"
)

// AdminHeaderUI creates the admin header navigation
func AdminHeaderUI(r *http.Request, homeURL string) hb.TagInterface {
	type navItem struct {
		title string
		href  string
		path  string
	}

	items := []navItem{
		{
			title: "Visitor Analytics",
			href:  UrlHome(r),
			path:  PathHome,
		},
		{
			title: "Visitor Activity",
			href:  UrlVisitorActivity(r),
			path:  PathVisitorActivity,
		},
		{
			title: "Visitor Paths",
			href:  UrlVisitorPaths(r),
			path:  PathVisitorPaths,
		},
		{
			title: "Page View Activity",
			href:  UrlPageViewActivity(r),
			path:  PathPageViewActivity,
		},
	}

	currentPath := strings.TrimSuffix(r.URL.Query().Get("path"), "/")
	if currentPath == "" {
		currentPath = strings.TrimSuffix(r.URL.Path, "/")
	}
	if currentPath == "" {
		currentPath = PathHome
	}

	nav := hb.Nav().
		Class("nav nav-pills nav-fill flex-column flex-sm-row gap-2").
		Attr("role", "tablist")

	for _, item := range items {
		itemPath := strings.TrimSuffix(item.path, "/")
		isActive := currentPath == itemPath

		linkClasses := "nav-link fw-semibold text-center"
		if isActive {
			linkClasses += " active"
		}

		link := hb.A().
			Class(linkClasses).
			Attr("role", "tab").
			Href(item.href).
			Text(item.title)

		if isActive {
			link = link.Attr("aria-current", "page")
		}

		nav = nav.Child(link)
	}

	return hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-center gap-3 mb-3").
		Child(hb.A().
			Class("navbar-brand fw-semibold text-decoration-none").
			Href(homeURL).
			HTML("Visitor Analytics")).
		Child(nav)
}

// CardUI creates a standard card component
func CardUI(title string, body hb.TagInterface) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML(title))).
		Child(hb.Div().
			Class("card-body").
			Child(body))
}

// StatCardUI creates a card displaying a single statistic
func StatCardUI(title string, value string, icon string, color string) hb.TagInterface {
	iconWrapperClasses := "rounded-circle d-flex align-items-center justify-content-center text-" + color
	backgroundClass := "bg-" + color + "-subtle"

	return hb.Div().
		Class("card h-100 shadow-sm border-0").
		Child(hb.Div().
			Class("card-body d-flex flex-column align-items-center justify-content-center gap-3 py-4").
			Child(hb.Div().
				Class(iconWrapperClasses + " " + backgroundClass).
				Style("width: 60px; height: 60px;").
				Child(hb.I().
					Class(icon + " fs-3"))).
			Child(hb.Small().
				Class("text-uppercase text-muted fw-semibold text-center letter-spacing-1").
				Text(title)).
			Child(hb.Heading3().
				Class("mb-0 fw-bold").
				Text(value)))
}

// NavCardUI creates a navigation card with icon and description
func NavCardUI(title string, href string, icon string, description string) hb.TagInterface {
	return hb.Div().
		Class("card h-100 shadow-sm").
		Style("transition: all 0.3s ease;").
		Child(hb.Div().
			Class("card-body text-center").
			Child(hb.Div().
				Class("mb-3").
				Child(hb.I().
					Class(icon + " fs-1"))).
			Child(hb.Heading5().
				Class("card-title").
				HTML(title)).
			Child(hb.P().
				Class("card-text text-muted mb-3").
				Text(description)).
			Child(hb.A().
				Class("btn btn-outline-primary").
				Href(href).
				Text("View Details")))
}

// PaginationUI creates a pagination component
func PaginationUI(currentPage int, totalPages int, urlFunc func(page int) string) hb.TagInterface {
	if totalPages <= 1 {
		return hb.Div()
	}

	nav := hb.Nav().
		Attr("aria-label", "Page navigation")

	ul := hb.UL().
		Class("pagination justify-content-center")

	// Previous button
	if currentPage > 1 {
		ul.Child(hb.LI().
			Class("page-item").
			Child(hb.A().
				Class("page-link").
				Href(urlFunc(currentPage - 1)).
				HTML("&laquo;")))
	} else {
		ul.Child(hb.LI().
			Class("page-item disabled").
			Child(hb.A().
				Class("page-link").
				Href("#").
				HTML("&laquo;")))
	}

	addPage := func(page int, active bool) {
		li := hb.LI().Class("page-item")
		link := hb.A().
			Class("page-link").
			Href(urlFunc(page)).
			Text(fmt.Sprintf("%d", page))

		if active {
			li.Class("page-item active")
			link = link.Attr("aria-current", "page")
		}

		li = li.Child(link)
		ul.Child(li)
	}

	addEllipsis := func() {
		ul.Child(hb.LI().
			Class("page-item disabled").
			Child(hb.Span().
				Class("page-link").
				HTML("&hellip;")))
	}

	if totalPages <= 7 {
		for i := 1; i <= totalPages; i++ {
			addPage(i, i == currentPage)
		}
	} else {
		// Always show the first page
		addPage(1, currentPage == 1)

		start := currentPage - 2
		if start < 2 {
			start = 2
		}

		end := currentPage + 2
		if end > totalPages-1 {
			end = totalPages - 1
		}

		if start > 2 {
			addEllipsis()
		}

		for i := start; i <= end; i++ {
			addPage(i, i == currentPage)
		}

		if end < totalPages-1 {
			addEllipsis()
		}

		addPage(totalPages, currentPage == totalPages)
	}

	// Next button
	if currentPage < totalPages {
		ul.Child(hb.LI().
			Class("page-item").
			Child(hb.A().
				Class("page-link").
				Href(urlFunc(currentPage + 1)).
				HTML("&raquo;")))
	} else {
		ul.Child(hb.LI().
			Class("page-item disabled").
			Child(hb.A().
				Class("page-link").
				Href("#").
				HTML("&raquo;")))
	}

	return nav.Child(ul)
}
