package add_salay_change

import (
	"salary_calculator/internal/dto/value_objects"
)

type In struct {
	Value float64                   `json:"value"`
	Date  *value_objects.SalaryDate `json:"date"`
}

type Out struct {
	Ok bool `json:"ok"`
}
