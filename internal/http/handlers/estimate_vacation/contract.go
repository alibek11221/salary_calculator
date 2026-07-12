package estimate_vacation

import (
	"context"

	"salary_calculator/internal/dto/estimate_vacation"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in estimate_vacation.In) (*estimate_vacation.Out, error)
}
