package list_bonuses

import (
	"context"

	"salary_calculator/internal/dto/list_bonuses"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context) (*list_bonuses.Out, error)
}
