package home

import "net/http"

// ControllerData contains the data needed for the home page
type ControllerData struct {
	Request      *http.Request
	dates        []string
	uniqueVisits []int64
	totalVisits  []int64
	firstVisits  []int64
	returnVisits []int64
	selectedPeriod string
	periodOptions  []periodOption
}

type periodOption struct {
	Value string
	Label string
}
