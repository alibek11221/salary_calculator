package value_objects

import (
	"salary_calculator/internal/services/work_days"
)

type SalaryCalculationContext struct {
	currentBase float64
	currentNDFL float64
	workdays    work_days.WorkdaysForMonth
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
