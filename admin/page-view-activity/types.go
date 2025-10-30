package pageviewactivity

import (
	"net/http"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

// ControllerData holds the data required to render the page view activity screen.
type ControllerData struct {
	Request    *http.Request
	Visitors   []statsstore.VisitorInterface
	Page       int
	TotalPages int
	PageSize   int
	TotalCount int64
	Filters    FilterOptions
}

// ControllerOptions aliases the shared controller options to avoid repetition in imports.
type ControllerOptions = shared.ControllerOptions

// FilterOptions represents the active filters on the page view activity screen.
type FilterOptions struct {
	Range   string
	From    string
	To      string
	Country string
	Device  string
	Browser string
}
