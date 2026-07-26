package home

import "testing"

import "github.com/dromara/carbon/v2"

func TestDatesInRange(t *testing.T) {
	start := carbon.CreateFromDate(2026, 7, 1, carbon.UTC)
	end := carbon.CreateFromDate(2026, 7, 3, carbon.UTC)

	dates := datesInRange(start, end)

	if len(dates) != 3 {
		t.Fatalf("expected 3 dates, got %d", len(dates))
	}

	if dates[0] != "2026-07-01" {
		t.Fatalf("expected first date 2026-07-01, got %s", dates[0])
	}

	if dates[2] != "2026-07-03" {
		t.Fatalf("expected last date 2026-07-03, got %s", dates[2])
	}
}

func TestDatesInRangeSingleDay(t *testing.T) {
	start := carbon.CreateFromDate(2026, 7, 1, carbon.UTC)
	end := carbon.CreateFromDate(2026, 7, 1, carbon.UTC)

	dates := datesInRange(start, end)

	if len(dates) != 1 {
		t.Fatalf("expected 1 date, got %d", len(dates))
	}
}

func TestDatesInRangeStartAfterEnd(t *testing.T) {
	start := carbon.CreateFromDate(2026, 7, 3, carbon.UTC)
	end := carbon.CreateFromDate(2026, 7, 1, carbon.UTC)

	dates := datesInRange(start, end)

	if len(dates) != 0 {
		t.Fatalf("expected 0 dates when start > end, got %d", len(dates))
	}
}
