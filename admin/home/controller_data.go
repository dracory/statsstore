package home

import (
	"net/http"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

// ControllerData contains the data needed for the home page
type ControllerData struct {
	Request        *http.Request
	visitors       []statsstore.VisitorInterface
	ui             shared.ControllerOptions
	dates          []string
	uniqueVisits   []int64
	totalVisits    []int64
	firstVisits    []int64
	returnVisits   []int64
	selectedPeriod string
	periodOptions  []periodOption

	liveVisitorCount        int64
	previousPeriodLabel     string
	previousPeriodUnique    int64
	previousPeriodTotal     int64
	previousPeriodFirst     int64
	previousPeriodReturning int64
	previousPeriodVisitors  []statsstore.VisitorInterface
	currentStats            extendedStats
	previousStats           extendedStats
}

type periodOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
