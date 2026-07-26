package home

import (
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

// ControllerData contains the data needed for traffic source computation.
type ControllerData struct {
	visitors []statsstore.VisitorInterface
	ui       shared.ControllerOptions
}

type periodOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
