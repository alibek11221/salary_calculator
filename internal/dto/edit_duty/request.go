package edit_duty

import (
	"salary_calculator/internal/dto/value_objects"
)

type In struct {
	Date       *value_objects.SalaryDate `json:"date"`
	InWorkdays int32                     `json:"in_workdays"`
	InHolidays int32                     `json:"in_holidays"`
}

type Out struct {
	Ok bool `json:"ok"`
}
