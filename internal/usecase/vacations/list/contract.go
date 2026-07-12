package list

import (
	"context"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	ListVacations(ctx context.Context) ([]dbstore.Vacation, error)
}
