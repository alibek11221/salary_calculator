package edit

import (
	"context"

	"salary_calculator/internal/generated/dbstore"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	GetDutyByDate(ctx context.Context, date string) (dbstore.Duty, error)
	UpdateDuty(ctx context.Context, arg dbstore.UpdateDutyParams) error
}
