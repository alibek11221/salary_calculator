package value_objects

import (
	"salary_calculator/internal/services/work_days"
)

type SalaryCalculationContext struct {
	currentBase      float64
	currentNDFL      float64
	workdays         work_days.WorkdaysForMonth
	vacationPayments []ExtraPayment
}

func NewSalaryContext(
	currentBase float64,
	currentNDFL float64,
	workdays work_days.WorkdaysForMonth,
) SalaryCalculationContext {
	return SalaryCalculationContext{
		currentBase: currentBase,
		currentNDFL: currentNDFL,
		workdays:    workdays,
	}
}

func (s *SalaryCalculationContext) CurrentBase() float64 {
	return s.currentBase
}

func (s *SalaryCalculationContext) CurrentNDFL() float64 {
	return s.currentNDFL
}

func (s *SalaryCalculationContext) Workdays() work_days.WorkdaysForMonth {
	return s.workdays
}

// WithVacationPayments возвращает копию контекста с net-выплатами отпускных.
// Суммы уже за вычетом НДФЛ, их готовит usecase.
func (s SalaryCalculationContext) WithVacationPayments(payments ...ExtraPayment) SalaryCalculationContext {
	s.vacationPayments = payments
	return s
}

func (s *SalaryCalculationContext) VacationPayments() []ExtraPayment {
	return s.vacationPayments
}
