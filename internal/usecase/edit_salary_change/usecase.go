package edit_salary_change

import (
	"context"
	"fmt"

	"salary_calculator/internal/dto/edit_salary_change"
	"salary_calculator/internal/generated/dbstore"

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

func (u *usecase) Do(ctx context.Context, in edit_salary_change.In) (*edit_salary_change.Out, error) {
	var id pgtype.UUID
	if err := id.Scan(in.ID); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	model := dbstore.UpdateChangeParams{
		ID:         id,
		Salary:     in.Value,
		ChangeFrom: in.Date.String(),
	}

	err := u.r.UpdateChange(ctx, model)
	if err != nil {
		return nil, err
	}

	return &edit_salary_change.Out{
		Ok: true,
	}, nil
}
