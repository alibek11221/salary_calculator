package get_salary_changes

import (
	"salary_calculator/internal/dto/value_objects"
)

type In struct{}

type Out struct {
	Changes []Change `json:"changes"`
}

type Change struct {
	ID    string                   `json:"id"`
	Value float64                  `json:"value"`
	Date  value_objects.SalaryDate `json:"date"`
}
