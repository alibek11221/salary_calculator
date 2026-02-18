package add

import (
	"context"
	"errors"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

var ErrDuplicateSalaryChange = errors.New("duplicate salary change")

type repo interface {
	InsertChange(ctx context.Context, arg dbstore.InsertChangeParams) error
}
