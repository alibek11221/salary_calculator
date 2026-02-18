package add

import (
	"context"
	"errors"

	"salary_calculator/internal/dto/add_duty"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/database"

	"github.com/jackc/pgx/v5/pgconn"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{
		r: r,
	}
}

func (u *usecase) Do(ctx context.Context, in add_duty.In) (*add_duty.Out, error) {
	if in.Date == nil {
		return nil, errors.New("date is required")
	}
	model := dbstore.InsertDutyParams{
		Date:       in.Date.String(),
		InWorkdays: in.InWorkdays,
		InHolidays: in.InHolidays,
	}

	err := u.r.InsertDuty(ctx, model)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.DuplicateEntryCode {
			return nil, ErrDuplicateDuty
		}

		return nil, err
	}

	return &add_duty.Out{
		Ok: true,
	}, nil
}
