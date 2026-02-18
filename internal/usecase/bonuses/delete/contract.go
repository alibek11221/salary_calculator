package delete

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	DeleteBonus(ctx context.Context, id pgtype.UUID) error
}
