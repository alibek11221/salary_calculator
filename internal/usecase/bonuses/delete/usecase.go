package delete

import (
	"context"

	"salary_calculator/internal/dto/delete_bonus"

	"github.com/jackc/pgx/v5/pgtype"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context, in delete_bonus.In) (*delete_bonus.Out, error) {
	var id pgtype.UUID
	if err := id.Scan(in.ID); err != nil {
		return nil, err
	}

	if err := u.r.DeleteBonus(ctx, id); err != nil {
		return nil, err
	}

	return &delete_bonus.Out{Ok: true}, nil
}
