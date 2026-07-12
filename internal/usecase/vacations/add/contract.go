package add

import (
	"context"
	"errors"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

var ErrOverlappingVacation = errors.New("vacation overlaps with existing one")

type repo interface {
	InsertVacation(ctx context.Context, arg dbstore.InsertVacationParams) error
}
