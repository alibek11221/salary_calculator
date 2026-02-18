package list_bonuses

import (
	"salary_calculator/internal/dto/value_objects"
)

type Out struct {
	Bonuses []Bonus `json:"bonuses"`
}

type Bonus struct {
	ID          string                    `json:"id"`
	Value       float64                   `json:"value"`
	Date        *value_objects.SalaryDate `json:"date"`
	Coefficient float64                   `json:"coefficient"`
}
