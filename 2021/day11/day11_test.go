package day11

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Load sample data for use?
	// Load result to do comparison with??
	os.Exit(m.Run())
}

func Test_hasAdjacentFlashed(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
		{"happy #1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasAdjacentFlashed(); got != tt.want {
				t.Errorf("hasAdjacentFlashed() = %v, want %v", got, tt.want)
			}
		})
	}
}
