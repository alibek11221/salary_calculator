package add

import (
	"context"
	"errors"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

var ErrDuplicateBonus = errors.New("duplicate bonus")

type repo interface {
	InsertBonus(ctx context.Context, arg dbstore.InsertBonusParams) error
}
