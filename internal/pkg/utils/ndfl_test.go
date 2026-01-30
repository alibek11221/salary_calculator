package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateNDFL(t *testing.T) {
	tests := []struct {
		name     string
		salary   float64
		expected float64
	}{
		{"Zero salary", 0, 0},
		{"Negative salary", -1000, 0},
		{"Low income 13%", 150000, 13.0},
		{"Boundary 13%", 200000, 13.0},
		{"15% bracket start", 200001, 15.0},
		{"15% bracket", 349000, 15.0},
		{"Boundary 15%", 416666.66, 15.0},
		{"18% bracket", 1000000, 18.0},
		{"20% bracket", 2000000, 20.0},
		{"22% bracket", 5000000, 22.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateNDFL(tt.salary)
			assert.InDelta(t, tt.expected, result, 0.01)
		})
	}
}
