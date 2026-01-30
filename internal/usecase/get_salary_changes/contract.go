package get_salary_changes

import (
	"context"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	ListChanges(ctx context.Context) ([]dbstore.SalaryChange, error)
}
