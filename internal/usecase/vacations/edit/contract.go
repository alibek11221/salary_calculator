package edit

import (
	"context"
	"errors"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

var (
	ErrOverlappingVacation = errors.New("vacation overlaps with existing one")
	ErrInvalidID           = errors.New("invalid vacation id")
)

type repo interface {
	UpdateVacation(ctx context.Context, arg dbstore.UpdateVacationParams) error
}
