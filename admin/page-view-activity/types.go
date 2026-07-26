package pageviewactivity

import (
	shared "github.com/dracory/statsstore/admin/shared"
)

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
