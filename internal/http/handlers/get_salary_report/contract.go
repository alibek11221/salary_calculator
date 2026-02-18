package get_salary_report

import (
	"context"

	"salary_calculator/internal/dto/get_salary_report"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in get_salary_report.In) (*get_salary_report.Out, error)
}
