package delete_duty

import (
	"context"

	"salary_calculator/internal/dto/delete_duty"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context, in delete_duty.In) (*delete_duty.Out, error)
}
