package delete_bonus

import (
	"context"

	"salary_calculator/internal/dto/delete_bonus"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in delete_bonus.In) (*delete_bonus.Out, error)
}
