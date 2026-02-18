package edit_duty

import (
	"context"

	"salary_calculator/internal/dto/edit_duty"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in edit_duty.In) (*edit_duty.Out, error)
}
