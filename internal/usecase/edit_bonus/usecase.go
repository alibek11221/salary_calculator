package edit_bonus

import (
	"context"
	"errors"
	"salary_calculator/internal/dto/edit_bonus"
	"salary_calculator/internal/generated/dbstore"

	"github.com/jackc/pgx/v5/pgtype"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context, in edit_bonus.In) (*edit_bonus.Out, error) {
	if in.Date == nil {
		return nil, errors.New("date is required")
	}
	var id pgtype.UUID
	if err := id.Scan(in.ID); err != nil {
		return nil, err
	}

	arg := dbstore.UpdateBonusParams{
		ID:          id,
		Value:       in.Value,
		Date:        in.Date.String(),
		Coefficient: in.Coefficient,
	}

	if err := u.r.UpdateBonus(ctx, arg); err != nil {
		return nil, err
	}

	return &edit_bonus.Out{Ok: true}, nil
}
