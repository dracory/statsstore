package shared

import (
	"fmt"

	"github.com/dracory/hb"
)

// FilterSelectOption represents an option in a select dropdown within the filter modal.
type FilterSelectOption struct {
	Value string
	Label string
}

// FilterFieldDef defines a single filter field rendered in the filter modal.
type FilterFieldDef struct {
	Name    string               // query param name (e.g. "range", "country")
	Label   string               // display label
	Type    string               // "select" or "text"
	Options []FilterSelectOption // only for Type == "select"
	Value   string               // current value (pre-populated)
}

// FilterModalConfig configures the filter modal.
type FilterModalConfig struct {
	ModalID string
	Title   string
	Fields  []FilterFieldDef
}

// FilterModalButton renders a button that opens the filter modal.
func FilterModalButton(modalID string) hb.TagInterface {
	return hb.Button().
		Class("btn btn-sm btn-outline-primary").
		Attr("type", "button").
		Attr("data-bs-toggle", "modal").
		Attr("data-bs-target", "#"+modalID).
		HTML(`<i class="bi bi-funnel"></i> Filters`)
}

// FilterModal renders a Bootstrap modal with filter form fields.
// On "Apply Filters", JavaScript collects non-empty field values, preserves
// per_page from the current URL, resets to page 1, and navigates.
// On "Clear Filters", all filter params are removed and the page reloads.
func FilterModal(config FilterModalConfig) hb.TagInterface {
	// Build form fields
	var fields []hb.TagInterface

	for _, field := range config.Fields {
		fieldID := config.ModalID + "-" + field.Name

		label := hb.Label().
			Class("form-label fw-semibold").
			Attr("for", fieldID).
			Text(field.Label)

		var input hb.TagInterface
		if field.Type == "select" {
			selectTag := hb.Select().
				Class("form-select").
				ID(fieldID).
				Attr("name", field.Name)

			// Add empty option
			selectTag = selectTag.Child(hb.Option().
				Attr("value", "").
				Text("(any)"))

			for _, opt := range field.Options {
				option := hb.Option().
					Attr("value", opt.Value).
					Text(opt.Label)
				if field.Value == opt.Value {
					option = option.Attr("selected", "selected")
				}
				selectTag = selectTag.Child(option)
			}
			input = selectTag
		} else {
			input = hb.Input().
				Class("form-control").
				ID(fieldID).
				Attr("type", "text").
				Attr("name", field.Name).
				Attr("value", field.Value).
				Attr("placeholder", "(any)")
		}

		fields = append(fields, hb.Div().
			Class("mb-3").
			Child(label).
			Child(input))
	}

	// Modal header
	modalTitle := hb.Heading5().
		Class("modal-title").
		ID(config.ModalID + "Label").
		Text(config.Title)

	closeButton := hb.Button().
		Class("btn-close").
		Attr("type", "button").
		Attr("data-bs-dismiss", "modal").
		Attr("aria-label", "Close")

	modalHeader := hb.Div().
		Class("modal-header").
		Child(modalTitle).
		Child(closeButton)

	// Modal body with form fields
	modalBody := hb.Div().
		Class("modal-body").
		Children(fields)

	// Modal footer with Clear + Apply buttons
	clearBtn := hb.Button().
		Class("btn btn-outline-secondary me-auto").
		Attr("type", "button").
		Attr("id", config.ModalID+"-clear").
		HTML(`<i class="bi bi-x-circle"></i> Clear All`)

	cancelBtn := hb.Button().
		Class("btn btn-outline-secondary").
		Attr("type", "button").
		Attr("data-bs-dismiss", "modal").
		Text("Cancel")

	applyBtn := hb.Button().
		Class("btn btn-primary").
		Attr("type", "button").
		Attr("id", config.ModalID+"-apply").
		HTML(`<i class="bi bi-check-lg"></i> Apply Filters`)

	modalFooter := hb.Div().
		Class("modal-footer").
		Child(clearBtn).
		Child(cancelBtn).
		Child(applyBtn)

	// Assemble modal
	modalContent := hb.Div().
		Class("modal-content").
		Child(modalHeader).
		Child(modalBody).
		Child(modalFooter)

	modalDialog := hb.Div().
		Class("modal-dialog modal-dialog-centered").
		Child(modalContent)

	modal := hb.Div().
		Class("modal fade").
		ID(config.ModalID).
		Attr("tabindex", "-1").
		Attr("aria-labelledby", config.ModalID+"Label").
		Attr("aria-hidden", "true").
		Child(modalDialog)

	// JavaScript for apply/clear behaviour
	script := hb.Script(fmt.Sprintf(`
		(function() {
			var modalId = '%s';

			function buildFilteredURL(clearAll) {
				var url = new URL(window.location.href);
				var preserveKeys = ['per_page', 'path'];
				var preserved = {};
				preserveKeys.forEach(function(key) {
					var val = url.searchParams.get(key);
					if (val !== null) { preserved[key] = val; }
				});
				url.search = '';
				if (!clearAll) {
					Object.keys(preserved).forEach(function(key) {
						url.searchParams.set(key, preserved[key]);
					});
				} else {
					if (preserved['path']) { url.searchParams.set('path', preserved['path']); }
				}
				if (!clearAll) {
					var modal = document.getElementById(modalId);
					if (modal) {
						modal.querySelectorAll('select[name], input[name]').forEach(function(field) {
							var val = (field.value || '').trim();
							if (val !== '') {
								url.searchParams.set(field.name, val);
							}
						});
					}
				}
				url.searchParams.set('page', '1');
				return url.toString();
			}

			var applyBtn = document.getElementById(modalId + '-apply');
			if (applyBtn) {
				applyBtn.addEventListener('click', function() {
					window.location.href = buildFilteredURL(false);
				});
			}

			var clearBtn = document.getElementById(modalId + '-clear');
			if (clearBtn) {
				clearBtn.addEventListener('click', function() {
					window.location.href = buildFilteredURL(true);
				});
			}
		})();
	`, config.ModalID))

	return hb.Div().Child(modal).Child(script)
}

// RangeFilterOptions returns the standard time-range select options.
func RangeFilterOptions() []FilterSelectOption {
	return []FilterSelectOption{
		{Value: "24h", Label: "Last 24 Hours"},
		{Value: "today", Label: "Today"},
		{Value: "7d", Label: "Last 7 Days"},
		{Value: "30d", Label: "Last 30 Days"},
	}
}

// DeviceFilterOptions returns the standard device select options.
func DeviceFilterOptions() []FilterSelectOption {
	return []FilterSelectOption{
		{Value: "desktop", Label: "Desktop"},
		{Value: "mobile", Label: "Mobile"},
		{Value: "tablet", Label: "Tablet"},
		{Value: "empty", Label: "Unknown"},
	}
}
