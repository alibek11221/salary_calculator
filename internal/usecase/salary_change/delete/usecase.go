package delete

import (
	"context"
	"fmt"

	"salary_calculator/internal/dto/delete_salary_change"

	"github.com/jackc/pgx/v5/pgtype"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{
		r: r,
	}
}

func (u *usecase) Do(ctx context.Context, in delete_salary_change.In) (*delete_salary_change.Out, error) {
	var id pgtype.UUID
	if err := id.Scan(in.ID); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	err := u.r.DeleteChange(ctx, id)
	if err != nil {
		return nil, err
	}

	return &delete_salary_change.Out{
		Ok: true,
	}, nil
}
