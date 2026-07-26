package home

import "testing"

func TestFormatSummaryDate(t *testing.T) {
	result := formatSummaryDate("2026-07-26")
	if result == "2026-07-26" {
		t.Fatalf("expected formatted date, got original: %s", result)
	}
}

func TestFormatSummaryDateInvalid(t *testing.T) {
	result := formatSummaryDate("not-a-date")
	if result != "not-a-date" {
		t.Fatalf("expected original string for invalid date, got: %s", result)
	}
}

func TestChangePercentInt(t *testing.T) {
	if result := changePercentInt(150, 100); result != 50.0 {
		t.Fatalf("expected 50.0, got %f", result)
	}
}

func TestChangePercentIntPreviousZero(t *testing.T) {
	if result := changePercentInt(100, 0); result != 100.0 {
		t.Fatalf("expected 100.0, got %f", result)
	}
}

func TestChangePercentIntBothZero(t *testing.T) {
	if result := changePercentInt(0, 0); result != 0.0 {
		t.Fatalf("expected 0.0, got %f", result)
	}
}

func TestChangePercentFloat(t *testing.T) {
	if result := changePercentFloat(150.0, 100.0); result != 50.0 {
		t.Fatalf("expected 50.0, got %f", result)
	}
}

func TestChangePercentFloatPreviousZero(t *testing.T) {
	if result := changePercentFloat(50.0, 0.0); result != 100.0 {
		t.Fatalf("expected 100.0, got %f", result)
	}
}

func TestChangePercentFloatBothZero(t *testing.T) {
	if result := changePercentFloat(0.0, 0.0); result != 0.0 {
		t.Fatalf("expected 0.0, got %f", result)
	}
}
