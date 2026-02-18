package delete

import (
	"context"
	"errors"

	"salary_calculator/internal/dto/delete_duty"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context, in delete_duty.In) (*delete_duty.Out, error) {
	if in.Date == nil {
		return nil, errors.New("date is required")
	}

	existingDuty, err := u.r.GetDutyByDate(ctx, in.Date.String())
	if err != nil {
		return nil, err
	}

	if err := u.r.DeleteDuty(ctx, existingDuty.ID); err != nil {
		return nil, err
	}

	return &delete_duty.Out{Ok: true}, nil
}
