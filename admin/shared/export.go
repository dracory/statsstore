package shared

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"
)

// ExportFilename generates a CSV filename with the current UTC date.
// Example: ExportFilename("visitor-paths") → "visitor-paths-2024-01-15.csv"
func ExportFilename(prefix string) string {
	return fmt.Sprintf("%s-%s.csv", prefix, time.Now().UTC().Format("2006-01-02"))
}

// ExportCSV writes CSV content to the response writer with proper headers and
// a UTF-8 BOM for Excel/Google Sheets compatibility. It returns the CSV string
// so callers (which return a string from their Handler) can pass it through.
//
// On error it sets a 500 status code and returns an error message string.
func ExportCSV(w http.ResponseWriter, filename string, headers []string, rows [][]string) string {
	buffer := &bytes.Buffer{}

	// UTF-8 BOM — ensures Excel and Google Sheets interpret encoding correctly
	buffer.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(buffer)

	if err := writer.Write(headers); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "Failed to generate export"
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return "Failed to generate export"
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "Failed to generate export"
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return buffer.String()
}
