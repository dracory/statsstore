package shared

import (
	"fmt"

	"github.com/gouniverse/hb"
)

// AdminHeaderUI creates the admin header navigation
func AdminHeaderUI(ui UIContext) hb.TagInterface {
	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(URL(ui.GetRequest(), ControllerHome, nil)).
		Class("nav-link")

	nav := hb.Nav().
		Class("navbar navbar-expand-lg navbar-light bg-light").
		Child(hb.Div().
			Class("container-fluid").
			Child(hb.A().
				Class("navbar-brand").
				Href(ui.GetHomeURL()).
				HTML("Kalleidoscope")).
			Child(hb.Div().
				Class("collapse navbar-collapse").
				Child(hb.Div().
					Class("navbar-nav").
					Child(linkHome))))

	return nav
}

// CardUI creates a standard card component
func CardUI(title string, body hb.TagInterface) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header bg-light").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML(title))).
		Child(hb.Div().
			Class("card-body").
			Child(body))
}

// StatCardUI creates a card displaying a single statistic
func StatCardUI(title string, value string, icon string, color string) hb.TagInterface {
	return hb.Div().
		Class("card h-100 border-0 shadow-sm").
		Child(hb.Div().
			Class("card-body").
			Child(hb.Div().
				Class("d-flex align-items-center").
				Child(hb.Div().
					Class("flex-shrink-0").
					Child(hb.I().
						Class(icon + " fs-1 text-" + color))).
				Child(hb.Div().
					Class("flex-grow-1 ms-3").
					Child(hb.P().
						Class("card-text text-muted mb-0").
						Text(title)).
					Child(hb.Heading3().
						Class("mb-0 fw-bold").
						Text(value)))))
}

// NavCardUI creates a navigation card with icon and description
func NavCardUI(title string, href string, icon string, description string) hb.TagInterface {
	return hb.Div().
		Class("card h-100 border-0 shadow-sm hover-shadow").
		Style("transition: all 0.3s ease;").
		Child(hb.Div().
			Class("card-body text-center").
			Child(hb.Div().
				Class("mb-3").
				Child(hb.I().
					Class(icon + " fs-1 text-primary"))).
			Child(hb.Heading5().
				Class("card-title").
				HTML(title)).
			Child(hb.P().
				Class("card-text text-muted").
				Text(description)).
			Child(hb.A().
				Class("btn btn-outline-primary mt-2").
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

	// Page numbers
	for i := 1; i <= totalPages; i++ {
		if i == currentPage {
			ul.Child(hb.LI().
				Class("page-item active").
				Child(hb.A().
					Class("page-link").
					Href(urlFunc(i)).
					Text(fmt.Sprintf("%d", i))))
		} else {
			ul.Child(hb.LI().
				Class("page-item").
				Child(hb.A().
					Class("page-link").
					Href(urlFunc(i)).
					Text(fmt.Sprintf("%d", i))))
		}
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
