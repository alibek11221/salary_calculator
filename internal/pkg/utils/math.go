package utils

import "math"

func SubPercentage(from float64, percent float64) float64 {
	return from * (1 - percent/100)
}

func ToTwoDecimals(num float64) float64 {
	return math.Round(num*100) / 100
}

func TenPercent(num float64) float64 {
	return num * 0.1
}
