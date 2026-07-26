package visitorpaths

import (
	"github.com/dracory/statsstore/admin/shared"
)

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
