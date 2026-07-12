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
		{"Low income fully in 13% bracket", 150000, 13.0},
		{"Exactly 2.4M annual is fully 13%", 200000, 13.0},
		// 4 188 000/год: 2 400 000×0.13 + 1 788 000×0.15 = 580 200 → 13.853868...%
		{"349000 monthly crosses into 15% bracket", 349000, 13.8538681948424},
		// 12 000 000/год: 312 000 + 390 000 + 7 000 000×0.18 = 1 962 000 → 16.35%
		{"18% bracket", 1000000, 16.35},
		// 24 000 000/год: 312 000 + 390 000 + 2 700 000 + 4 000 000×0.20 = 4 202 000 → 17.508333...%
		{"20% bracket", 2000000, 17.5083333333333},
		// 60 000 000/год: 312 000 + 390 000 + 2 700 000 + 6 000 000 + 10 000 000×0.22 = 11 602 000 → 19.336666...%
		{"22% bracket", 5000000, 19.3366666666667},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateNDFL(tt.salary)
			assert.InDelta(t, tt.expected, result, 1e-6)
		})
	}

	// Защита от off-by-one на границе первой ступени: чуть выше 2.4М/год
	// эффективная ставка должна стать строго больше 13%.
	t.Run("Just above 2.4M annual is strictly greater than 13", func(t *testing.T) {
		assert.Greater(t, CalculateNDFL(200000.01), 13.0)
	})
}
