package get_salary_report

import (
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/vacation_pay"
)

type In struct {
	Year  int `json:"year"`
	Month int `json:"month"`

	// ExtraVacations — гипотетические отпуска, не сохранённые в БД.
	// Через HTTP не передаются, заполняет usecase vacations/estimate.
	ExtraVacations []dbstore.Vacation `json:"-"`

	// Earnings — предзагруженные данные для среднего заработка.
	// Через HTTP не передаются, заполняет usecase vacations/estimate;
	// если non-nil, usecase отчёта не грузит их из БД повторно.
	Earnings *vacation_pay.EarningsData `json:"-"`
}

type Out struct {
	BaseSalary float64                             `json:"base_salary"`
	Result     *calculator.SalaryCalculationResult `json:"result,omitempty"`
}
