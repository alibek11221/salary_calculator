package list_vacations

import (
	"context"

	"salary_calculator/internal/dto/list_vacations"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type usecase interface {
	Do(ctx context.Context) (*list_vacations.Out, error)
}
