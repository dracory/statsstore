package visitoractivity

import (
	shared "github.com/dracory/statsstore/admin/shared"
)

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
