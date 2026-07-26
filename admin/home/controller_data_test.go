package home

import "testing"

func TestControllerDataZeroValue(t *testing.T) {
	var data ControllerData
	if data.selectedPeriod != "" {
		t.Fatalf("expected empty selectedPeriod, got %s", data.selectedPeriod)
	}
	if data.liveVisitorCount != 0 {
		t.Fatalf("expected zero liveVisitorCount, got %d", data.liveVisitorCount)
	}
}

func TestPeriodOptionFields(t *testing.T) {
	opt := periodOption{Value: "today", Label: "Today"}
	if opt.Value != "today" || opt.Label != "Today" {
		t.Fatalf("unexpected periodOption: %+v", opt)
	}
}
