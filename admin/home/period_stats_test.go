package home

import "testing"

import "github.com/dracory/statsstore"

import "github.com/dromara/carbon/v2"

import "time"

func TestComputePeriodStatsEmpty(t *testing.T) {
	stats := computePeriodStats(nil, []string{})
	if len(stats.dates) != 0 {
		t.Fatalf("expected 0 dates, got %d", len(stats.dates))
	}
}

func TestComputePeriodStatsNoDates(t *testing.T) {
	base := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	v := statsstore.NewVisitor().
		SetIpAddress("1.1.1.1").
		SetCreatedAt(carbon.CreateFromStdTime(base).ToDateTimeString())

	stats := computePeriodStats([]statsstore.VisitorInterface{v}, []string{})
	if stats.totalTotal != 0 {
		t.Fatalf("expected 0 total, got %d", stats.totalTotal)
	}
}

func TestComputePeriodStatsUnknownIP(t *testing.T) {
	base := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	v := statsstore.NewVisitor().
		SetCreatedAt(carbon.CreateFromStdTime(base).ToDateTimeString())

	stats := computePeriodStats([]statsstore.VisitorInterface{v}, []string{"2026-07-13"})
	if stats.totalUnique != 1 {
		t.Fatalf("expected 1 unique (unknown-ip fallback), got %d", stats.totalUnique)
	}
}
