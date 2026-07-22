package utils

import "math"

// ndflBracket — ступень прогрессивной шкалы НДФЛ: часть годового дохода
// до границы upTo облагается по ставке rate.
type ndflBracket struct {
	upTo float64
	rate float64
}

// Прогрессивная шкала НДФЛ с 2025 года (п. 1 ст. 224 НК РФ).
var ndflBrackets = []ndflBracket{
	{upTo: 2_400_000, rate: 0.13},
	{upTo: 5_000_000, rate: 0.15},
	{upTo: 20_000_000, rate: 0.18},
	{upTo: 50_000_000, rate: 0.20},
	{upTo: math.Inf(1), rate: 0.22},
}

// CalculateNDFL возвращает эффективную ставку НДФЛ в процентах для месячного
// оклада salary: налог считается ступенчато от годового дохода (salary × 12),
// возвращается отношение суммарного налога к годовому доходу.
// Для salary <= 0 возвращается 0.
func CalculateNDFL(salary float64) float64 {
	if salary <= 0 {
		return 0
	}

	annual := salary * 12

	tax, prevUpTo := 0.0, 0.0
	for _, b := range ndflBrackets {
		if annual <= b.upTo {
			tax += (annual - prevUpTo) * b.rate
			break
		}
		tax += (b.upTo - prevUpTo) * b.rate
		prevUpTo = b.upTo
	}

	return tax / annual * 100
}
