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

	modalBody := hb.Div().
		Class("modal-body").
		ID("visitorDetailModalContent")

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
