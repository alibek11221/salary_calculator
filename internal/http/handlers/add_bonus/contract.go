package add_bonus

import (
	"context"

	"salary_calculator/internal/dto/add_bonus"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in add_bonus.In) (*add_bonus.Out, error)
}
