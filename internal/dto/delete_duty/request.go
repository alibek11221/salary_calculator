package delete_duty

import (
	"salary_calculator/internal/dto/value_objects"
)

type In struct {
	Date *value_objects.SalaryDate `json:"date"`
}

type Out struct {
	Ok bool `json:"ok"`
}
