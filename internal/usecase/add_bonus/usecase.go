package add_bonus

import (
	"context"
	"errors"
	"salary_calculator/internal/dto/add_bonus"
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

func (u *usecase) Do(ctx context.Context, in add_bonus.In) (*add_bonus.Out, error) {
	if in.Date == nil {
		return nil, errors.New("date is required")
	}
	model := dbstore.InsertBonusParams{
		Value:       in.Value,
		Date:        in.Date.String(),
		Coefficient: in.Coefficient,
	}

	err := u.r.InsertBonus(ctx, model)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.DuplicateEntryCode {
			return nil, ErrDuplicateBonus
		}

		return nil, err
	}

	return &add_bonus.Out{
		Ok: true,
	}, nil
}
