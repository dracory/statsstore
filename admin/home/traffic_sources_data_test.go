package home

import (
	"testing"
	"time"

	"github.com/dracory/statsstore"
	"github.com/dromara/carbon/v2"
)

func newTestVisitor(fingerprint, ip string, t time.Time) statsstore.VisitorInterface {
	v := statsstore.NewVisitor()
	v.SetFingerprint(fingerprint)
	v.SetIpAddress(ip)
	v.SetCreatedAt(carbon.CreateFromStdTime(t).ToDateTimeString())
	return v
}

func TestComputeStatsOverview(t *testing.T) {
	base := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	visitors := []statsstore.VisitorInterface{
		newTestVisitor("fp1", "1.1.1.1", base),
		newTestVisitor("fp1", "1.1.1.1", base.Add(30*time.Second)),
		newTestVisitor("fp2", "2.2.2.2", base),
		newTestVisitor("fp3", "3.3.3.3", base.Add(2*time.Minute)),
	}

	stats := computeStatsOverview(visitors)

	if stats.BounceRate != "66.7%" {
		t.Errorf("bounce rate: expected 66.7%%, got %s", stats.BounceRate)
	}
	if stats.SessionDuration != "30s" {
		t.Errorf("avg visit duration: expected 30s, got %s", stats.SessionDuration)
	}
	if stats.Pageviews != "4" {
		t.Errorf("pageviews: expected 4, got %s", stats.Pageviews)
	}
	if stats.Sessions != "3" {
		t.Errorf("sessions: expected 3, got %s", stats.Sessions)
	}
}

func TestComputeStatsOverviewUsesFingerprintCalculateFallback(t *testing.T) {
	base := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	v1 := statsstore.NewVisitor().SetIpAddress("10.0.0.1").SetUserAgent("ua1")
	v1.SetCreatedAt(carbon.CreateFromStdTime(base).ToDateTimeString())
	v2 := statsstore.NewVisitor().SetIpAddress("10.0.0.1").SetUserAgent("ua1")
	v2.SetCreatedAt(carbon.CreateFromStdTime(base.Add(10 * time.Second)).ToDateTimeString())
	v3 := statsstore.NewVisitor().SetIpAddress("10.0.0.2").SetUserAgent("ua2")
	v3.SetCreatedAt(carbon.CreateFromStdTime(base).ToDateTimeString())

	stats := computeStatsOverview([]statsstore.VisitorInterface{v1, v2, v3})

	if stats.Sessions != "2" {
		t.Errorf("sessions: expected 2, got %s", stats.Sessions)
	}
	if stats.BounceRate != "50.0%" {
		t.Errorf("bounce rate: expected 50.0%%, got %s", stats.BounceRate)
	}
}

func TestComputePeriodStats(t *testing.T) {
	baseTime := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	base := carbon.CreateFromStdTime(baseTime)
	next := carbon.CreateFromStdTime(baseTime.AddDate(0, 0, 1))

	visitors := []statsstore.VisitorInterface{
		newTestVisitor("a", "1.1.1.1", baseTime),
		newTestVisitor("b", "2.2.2.2", baseTime.AddDate(0, 0, 1)),
		newTestVisitor("c", "1.1.1.1", baseTime.AddDate(0, 0, 1)),
	}

	stats := computePeriodStats(visitors, []string{base.ToDateString(), next.ToDateString()})

	if stats.totalUnique != 3 {
		t.Errorf("total unique: expected 3, got %d", stats.totalUnique)
	}
	if stats.totalTotal != 3 {
		t.Errorf("total total: expected 3, got %d", stats.totalTotal)
	}
	if stats.totalFirst != 2 {
		t.Errorf("total first: expected 2, got %d", stats.totalFirst)
	}
	if stats.totalReturning != 1 {
		t.Errorf("total returning: expected 1, got %d", stats.totalReturning)
	}
	if len(stats.dates) != 2 || stats.dates[0] != base.ToDateString() {
		t.Errorf("dates: expected %v, got %v", []string{base.ToDateString(), next.ToDateString()}, stats.dates)
	}
}
