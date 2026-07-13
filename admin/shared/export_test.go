package shared

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExportFilename(t *testing.T) {
	filename := ExportFilename("visitor-paths")
	if !strings.HasPrefix(filename, "visitor-paths-") {
		t.Fatalf("expected filename to start with 'visitor-paths-', got: %s", filename)
	}
	if !strings.HasSuffix(filename, ".csv") {
		t.Fatalf("expected filename to end with '.csv', got: %s", filename)
	}
}

func TestExportCSV(t *testing.T) {
	headers := []string{"Name", "Country", "IP"}
	rows := [][]string{
		{"Alice", "US", "127.0.0.1"},
		{"Bob", "DE", "192.168.0.1"},
	}

	rr := httptest.NewRecorder()
	result := ExportCSV(rr, "test-export.csv", headers, rows)

	if status := rr.Result().StatusCode; status != http.StatusOK {
		t.Fatalf("unexpected status: %d", status)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "text/csv; charset=utf-8" {
		t.Fatalf("unexpected content type: %s", ct)
	}

	cd := rr.Header().Get("Content-Disposition")
	if cd != "attachment; filename=\"test-export.csv\"" {
		t.Fatalf("unexpected content disposition: %s", cd)
	}

	if !strings.HasPrefix(result, "\xEF\xBB\xBF") {
		t.Fatalf("expected UTF-8 BOM at start of result")
	}

	records, err := csv.NewReader(strings.NewReader(strings.TrimPrefix(result, "\xEF\xBB\xBF"))).ReadAll()
	if err != nil {
		t.Fatalf("failed to parse csv: %v", err)
	}

	if len(records) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(records))
	}

	expectedHeader := []string{"Name", "Country", "IP"}
	for i, v := range records[0] {
		if v != expectedHeader[i] {
			t.Fatalf("unexpected header[%d]: got %q, want %q", i, v, expectedHeader[i])
		}
	}

	if records[1][0] != "Alice" || records[2][0] != "Bob" {
		t.Fatalf("unexpected data rows: %+v", records[1:])
	}
}

func TestExportCSVBOMStrippedFromFirstField(t *testing.T) {
	headers := []string{"Col1"}
	rows := [][]string{{"val1"}}

	rr := httptest.NewRecorder()
	body := ExportCSV(rr, "test.csv", headers, rows)

	// The BOM should be present in the raw body
	if !strings.HasPrefix(body, "\xEF\xBB\xBF") {
		t.Fatalf("expected BOM prefix")
	}

	// After stripping BOM, csv parsing should yield clean header
	stripped := strings.TrimPrefix(body, "\xEF\xBB\xBF")
	records, err := csv.NewReader(strings.NewReader(stripped)).ReadAll()
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	if records[0][0] != "Col1" {
		t.Fatalf("expected 'Col1', got %q", records[0][0])
	}
}
