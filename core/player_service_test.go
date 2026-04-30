package core

import (
	"testing"
)

func TestCalculateTransition(t *testing.T) {
	tests := []struct {
		name       string
		totalTime  int
		mediaCount int
		expected   int
	}{
		{"basic", 180, 10, 18},
		{"one item", 60, 1, 60},
		{"zero media", 180, 0, 0},
		{"negative media", 180, -1, 0},
		{"large set", 3600, 100, 36},
		{"short time", 5, 3, 1},
		{"zero time", 0, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTransition(tt.totalTime, tt.mediaCount)
			if result != tt.expected {
				t.Errorf("CalculateTransition(%d, %d) = %d; want %d",
					tt.totalTime, tt.mediaCount, result, tt.expected)
			}
		})
	}
}
