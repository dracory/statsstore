package settings

import (
	"net/http"

	shared "github.com/dracory/statsstore/admin/shared"
)

// ControllerOptions configures the settings controller views
type ControllerOptions = shared.ControllerOptions

// ControllerData contains the data needed for the settings page
type ControllerData struct {
	Request      *http.Request
	ExcludedIPs  []string
	ErrorMessage string
}
