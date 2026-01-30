package get_bonuses

import (
	"context"

	"salary_calculator/internal/dto/get_bonuses"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context) (*get_bonuses.Out, error)
}
