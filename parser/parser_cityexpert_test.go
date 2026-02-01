package main

import (
	"testing"
)

func TestCityExpertPaginationParsing(t *testing.T) {
	// <span ...>571-596 od 596 rezultata </span>
	textContent := "571-596 od 596 rezultata"

	total := parseCityExpertTotalCount(textContent)
	if total != 596 {
		t.Errorf("Expected 596, got %d", total)
	}

	// Test case with garbage or missing parts
	zero := parseCityExpertTotalCount("Some random text")
	if zero != 0 {
		t.Errorf("Expected 0 for invalid text, got %d", zero)
	}
}
