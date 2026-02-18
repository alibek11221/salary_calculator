package delete_s_change

import (
	"context"

	"salary_calculator/internal/dto/delete_salary_change"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in delete_salary_change.In) (*delete_salary_change.Out, error)
}
