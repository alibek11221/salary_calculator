package add_duty

import (
	"context"

	"salary_calculator/internal/dto/add_duty"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in add_duty.In) (*add_duty.Out, error)
}
