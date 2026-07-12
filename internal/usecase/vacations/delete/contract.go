package delete

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

var ErrInvalidID = errors.New("invalid vacation id")

type repo interface {
	DeleteVacation(ctx context.Context, id pgtype.UUID) error
}
