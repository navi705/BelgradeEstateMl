package main

import "testing"

func TestGetFloorLabel(t *testing.T) {
	tests := []struct {
		input    float32
		expected string
	}{
		{0.0, "Prizemlje"},
		{-3.0, "Suteren"},
		{0.5, "Visokoprizemlje"},
		{1000.0, "Potkrovlje"},
		{1.0, "1"},
		{2.0, "2"},
		{-5.0, "Nepoznato"},
	}

	for _, tt := range tests {
		got := GetFloorLabel(tt.input)
		if got != tt.expected {
			t.Errorf("GetFloorLabel(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestNormalizeFloorValue(t *testing.T) {
	tests := []struct {
		input    float32
		total    float32
		expected float32
	}{
		{1000.0, 5.0, 5.0},
		{1000.0, 10.0, 10.0},
		{1000.0, -5.0, 1.0},
		{0.0, 5.0, 0.0},
		{-3.0, 5.0, -3.0},
		{-5.0, 5.0, 0.0},
	}

	for _, tt := range tests {
		got := NormalizeFloorValue(tt.input, tt.total)
		if got != tt.expected {
			t.Errorf("NormalizeFloorValue(%v, %v) = %v, want %v", tt.input, tt.total, got, tt.expected)
		}
	}
}
