package get_salary_changes

import (
	"context"

	"salary_calculator/internal/dto/get_salary_changes"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context) (*get_salary_changes.Out, error)
}
