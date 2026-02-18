package list_duties

import (
	"context"

	"salary_calculator/internal/dto/list_duties"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context) (*list_duties.Out, error)
}
