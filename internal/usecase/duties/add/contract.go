package add

import (
	"context"
	"errors"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

var ErrDuplicateDuty = errors.New("duplicate duty")

type repo interface {
	InsertDuty(ctx context.Context, arg dbstore.InsertDutyParams) error
}
