package visitoractivity

import (
	"net/http"

	statsstore "github.com/dracory/statsstore"
	shared "github.com/dracory/statsstore/admin/shared"
)

// ControllerData contains the data needed for visitor activity pages
type ControllerData struct {
	Request    *http.Request
	Visitors   []statsstore.VisitorInterface
	Page       int
	TotalPages int
	PageSize   int
	TotalCount int64
	Filters    FilterOptions
}

// ControllerOptions configures the visitor activity controller views
type ControllerOptions = shared.ControllerOptions

// FilterOptions describes the active filters applied to the visitor list
type FilterOptions struct {
	Range   string
	From    string
	To      string
	Country string
	Device  string
}
