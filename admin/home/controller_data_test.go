package home

import "testing"

func TestControllerDataZeroValue(t *testing.T) {
	var data ControllerData
	if data.visitors != nil {
		t.Fatalf("expected nil visitors, got %v", data.visitors)
	}
}

func TestPeriodOptionFields(t *testing.T) {
	opt := periodOption{Value: "today", Label: "Today"}
	if opt.Value != "today" || opt.Label != "Today" {
		t.Fatalf("unexpected periodOption: %+v", opt)
	}
}
