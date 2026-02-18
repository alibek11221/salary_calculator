package list

import (
	"context"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	ListDuties(ctx context.Context) ([]dbstore.Duty, error)
}
