package calculator

import (
	"context"

	"salary_calculator/internal/generated/dbstore"
)

type repo interface {
	GetDutyByDate(ctx context.Context, date string) (dbstore.Duty, error)
	GetBonusByDate(ctx context.Context, date string) (dbstore.Bonuse, error)
}
