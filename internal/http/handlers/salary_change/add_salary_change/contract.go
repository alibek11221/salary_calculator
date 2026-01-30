package add_salary_change

import (
	"context"

	"salary_calculator/internal/dto/add_salay_change"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in add_salay_change.In) (*add_salay_change.Out, error)
}
