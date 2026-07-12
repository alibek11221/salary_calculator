package delete_vacation

import (
	"context"

	"salary_calculator/internal/dto/delete_vacation"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in delete_vacation.In) (*delete_vacation.Out, error)
}
