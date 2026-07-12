package add

import (
	"context"
	"errors"

	"salary_calculator/internal/dto/add_vacation"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/database"
	"salary_calculator/internal/services/vacation_pay"

	"github.com/jackc/pgx/v5/pgconn"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context, in add_vacation.In) (*add_vacation.Out, error) {
	if err := vacation_pay.ValidatePeriod(in.From.Time, in.To.Time); err != nil {
		return nil, err
	}

	err := u.r.InsertVacation(ctx, dbstore.InsertVacationParams{
		DateFrom: in.From.PG(),
		DateTo:   in.To.PG(),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ExclusionViolationCode {
			return nil, ErrOverlappingVacation
		}

		return nil, err
	}

	return &add_vacation.Out{Ok: true}, nil
}
