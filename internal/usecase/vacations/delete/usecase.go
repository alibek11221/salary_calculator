package delete

import (
	"context"

	"salary_calculator/internal/dto/delete_vacation"

	"github.com/jackc/pgx/v5/pgtype"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context, in delete_vacation.In) (*delete_vacation.Out, error) {
	var id pgtype.UUID
	if err := id.Scan(in.ID); err != nil {
		return nil, ErrInvalidID
	}

	if err := u.r.DeleteVacation(ctx, id); err != nil {
		return nil, err
	}

	return &delete_vacation.Out{Ok: true}, nil
}
