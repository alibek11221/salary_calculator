package edit_s_change

import (
	"context"

	"salary_calculator/internal/dto/edit_salary_change"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in edit_salary_change.In) (*edit_salary_change.Out, error)
}
