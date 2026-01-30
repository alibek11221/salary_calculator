package edit_bonus

import (
	"salary_calculator/internal/dto/value_objects"
)

type In struct {
	ID          string                    `json:"id"`
	Value       float64                   `json:"value"`
	Date        *value_objects.SalaryDate `json:"date"`
	Coefficient float64                   `json:"coefficient"`
}

type Out struct {
	Ok bool `json:"ok"`
}
