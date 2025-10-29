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
}

// ControllerOptions configures the visitor activity controller views
type ControllerOptions = shared.ControllerOptions
