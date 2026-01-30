package get_bonuses

import (
	"context"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	ListBonuses(ctx context.Context) ([]dbstore.Bonuse, error)
}
