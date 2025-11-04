package visitorpaths

import (
	"net/http"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

// ControllerData contains the data needed for visitor paths views
type visitorPathsControllerData struct {
	Request    *http.Request
	Paths      []statsstore.VisitorInterface
	Page       int
	TotalPages int
	PageSize   int
	TotalCount int64
	Filters    FilterOptions
}

// ControllerOptions alias for shared controller options
type ControllerOptions = shared.ControllerOptions

// FilterOptions captures the active filters for visitor paths
type FilterOptions struct {
	Range        string
	From         string
	To           string
	Country      string
	PathContains string
	PathExact    string
	Device       string
}
