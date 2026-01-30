package add_salary_change

import (
	"context"
	"errors"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

const DuplicateEntryCode = "23505"

var ErrDuplicateSalaryChange = errors.New("duplicate salary change")

type repo interface {
	InsertChange(ctx context.Context, arg dbstore.InsertChangeParams) error
}
