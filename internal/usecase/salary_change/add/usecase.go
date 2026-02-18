package add

import (
	"context"
	"errors"

	"salary_calculator/internal/dto/add_salary_change"
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

func (u *usecase) Do(ctx context.Context, in add_salary_change.In) (*add_salary_change.Out, error) {
	if in.Date == nil {
		return nil, errors.New("date is required")
	}
	model := dbstore.InsertChangeParams{
		Salary:     in.Value,
		ChangeFrom: in.Date.String(),
	}

	err := u.r.InsertChange(ctx, model)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.DuplicateEntryCode {
			return nil, ErrDuplicateSalaryChange
		}

		return nil, err
	}

	return &add_salary_change.Out{
		Ok: true,
	}, nil
}
