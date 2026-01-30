package get_salary_report

import (
	"salary_calculator/internal/services/calculator"
)

type In struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}
type Out struct {
	BaseSalary float64                            `json:"base_salary"`
	Result     calculator.SalaryCalculationResult `json:"result"`
}
