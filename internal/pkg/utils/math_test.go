package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubPercentage(t *testing.T) {
	tests := []struct {
		name     string
		from     float64
		percent  float64
		expected float64
	}{
		{"10% from 100", 100, 10, 90},
		{"0% from 100", 100, 0, 100},
		{"100% from 100", 100, 100, 0},
		{"13% from 1000", 1000, 13, 870},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SubPercentage(tt.from, tt.percent))
		})
	}
}

func TestToTwoDecimals(t *testing.T) {
	tests := []struct {
		name     string
		num      float64
		expected float64
	}{
		{"round up", 10.126, 10.13},
		{"round down", 10.124, 10.12},
		{"no round", 10.12, 10.12},
		{"one decimal", 10.1, 10.1},
		{"integer", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToTwoDecimals(tt.num))
		})
	}
}

func TestTenPercent(t *testing.T) {
	tests := []struct {
		name     string
		num      float64
		expected float64
	}{
		{"10% of 100", 100, 10},
		{"10% of 1000", 1000, 100},
		{"10% of 0", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TenPercent(tt.num))
		})
	}
}
