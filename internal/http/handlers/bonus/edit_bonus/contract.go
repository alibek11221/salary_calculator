package edit_bonus

import (
	"context"

	"salary_calculator/internal/dto/edit_bonus"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in edit_bonus.In) (*edit_bonus.Out, error)
}
