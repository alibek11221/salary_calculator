package delete

import (
	"context"

	"salary_calculator/internal/generated/dbstore"

	"github.com/jackc/pgx/v5/pgtype"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	GetDutyByDate(ctx context.Context, date string) (dbstore.Duty, error)
	DeleteDuty(ctx context.Context, id pgtype.UUID) error
}
