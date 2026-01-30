package add_salary_change

import (
	"context"
	"errors"

	"salary_calculator/internal/dto/add_salay_change"
	"salary_calculator/internal/generated/dbstore"

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

func (u *usecase) Do(ctx context.Context, in add_salay_change.In) (*add_salay_change.Out, error) {
	model := dbstore.InsertChangeParams{
		Salary:     in.Value,
		ChangeFrom: in.Date.String(),
	}

	err := u.r.InsertChange(ctx, model)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == DuplicateEntryCode {
			return nil, ErrDuplicateSalaryChange
		}

		return nil, err
	}

	return &add_salay_change.Out{
		Ok: true,
	}, nil
}
