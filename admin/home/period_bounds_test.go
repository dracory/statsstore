package home

import "testing"

func TestMaxInt(t *testing.T) {
	if maxInt(3, 5) != 5 {
		t.Fatalf("expected 5, got %d", maxInt(3, 5))
	}
	if maxInt(10, 2) != 10 {
		t.Fatalf("expected 10, got %d", maxInt(10, 2))
	}
	if maxInt(7, 7) != 7 {
		t.Fatalf("expected 7, got %d", maxInt(7, 7))
	}
}
