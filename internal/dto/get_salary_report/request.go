package get_salary_report

import (
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/services/calculator"
)

type In struct {
	Year  int `json:"year"`
	Month int `json:"month"`

	// ExtraVacations — гипотетические отпуска, не сохранённые в БД.
	// Через HTTP не передаются, заполняет usecase vacations/estimate.
	ExtraVacations []dbstore.Vacation `json:"-"`
}

type Out struct {
	BaseSalary float64                             `json:"base_salary"`
	Result     *calculator.SalaryCalculationResult `json:"result,omitempty"`
}
