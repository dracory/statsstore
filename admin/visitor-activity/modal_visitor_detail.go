package visitoractivity

import (
	"github.com/dracory/hb"
)

func visitorDetailModal() hb.TagInterface {
	modalTitle := hb.Heading5().
		Class("modal-title").
		ID("visitorDetailModalLabel").
		Text("Visitor Details")

	closeButton := hb.Button().
		Class("btn-close").
		Attr("type", "button").
		Attr("data-bs-dismiss", "modal").
		Attr("aria-label", "Close")

	modalHeader := hb.Div().
		Class("modal-header").
		Child(modalTitle).
		Child(closeButton)

	loadingSpinner := hb.Div().
		Class("spinner-border text-primary").
		Attr("role", "status").
		Child(hb.Span().Class("visually-hidden").Text("Loading..."))

	initialBody := hb.Div().
		Class("text-center p-4").
		Child(loadingSpinner)

	modalBody := hb.Div().
		Class("modal-body").
		ID("visitorDetailModalContent").
		Child(initialBody)

	modalFooter := hb.Div().
		Class("modal-footer").
		Child(hb.Button().
			Class("btn btn-secondary").
			Attr("type", "button").
			Attr("data-bs-dismiss", "modal").
			Text("Close"))

	modalContent := hb.Div().
		Class("modal-content").
		Child(modalHeader).
		Child(modalBody).
		Child(modalFooter)

	modalDialog := hb.Div().
		Class("modal-dialog modal-lg modal-dialog-scrollable").
		Child(modalContent)

	return hb.Div().
		Class("modal fade").
		ID("visitorDetailModal").
		Attr("tabindex", "-1").
		Attr("aria-labelledby", "visitorDetailModalLabel").
		Attr("aria-hidden", "true").
		Child(modalDialog)
}
