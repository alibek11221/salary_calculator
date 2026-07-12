package add_vacation

import (
	"context"

	"salary_calculator/internal/dto/add_vacation"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in add_vacation.In) (*add_vacation.Out, error)
}
