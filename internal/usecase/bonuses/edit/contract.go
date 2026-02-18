package edit

import (
	"context"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	UpdateBonus(ctx context.Context, arg dbstore.UpdateBonusParams) error
}
