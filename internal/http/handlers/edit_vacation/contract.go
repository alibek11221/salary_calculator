package edit_vacation

import (
	"context"

	"salary_calculator/internal/dto/edit_vacation"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in edit_vacation.In) (*edit_vacation.Out, error)
}
