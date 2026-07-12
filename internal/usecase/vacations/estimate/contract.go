package estimate

import (
	"context"
	"time"

	report_dto "salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/services/vacation_pay"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	GetLatestChangeBeforeDate(ctx context.Context, changeFrom string) (dbstore.SalaryChange, error)
}

type vacationPayService interface {
	LoadEarningsData(ctx context.Context) (*vacation_pay.EarningsData, error)
	CalculatePayWith(data *vacation_pay.EarningsData, from, to time.Time) (*vacation_pay.Pay, error)
}

type reportUsecase interface {
	Do(ctx context.Context, in report_dto.In) (*report_dto.Out, error)
}
