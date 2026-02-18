package get_salary_report

import (
	"context"

	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/http/work_calendar"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/work_days"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	GetLatestChangeBeforeDate(ctx context.Context, changeFrom string) (dbstore.SalaryChange, error)
}

type salaryCalculator interface {
	CalculateSalary(
		ctx context.Context,
		date value_objects.SalaryDate,
		sCtx value_objects.SalaryCalculationContext,
	) (*calculator.SalaryCalculationResult, error)
}

type workdaysClient interface {
	GetWorkdaysForMonth(ctx context.Context, month int, year int) (*work_calendar.WorkdayResponse, error)
}

type workdaysCalculator interface {
	CalculateWorkDaysForMonth(month *work_calendar.WorkdayResponse) *work_days.WorkdaysForMonth
}
