package edit

import (
	"context"
	"errors"

	"salary_calculator/internal/dto/edit_duty"
	"salary_calculator/internal/generated/dbstore"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context, in edit_duty.In) (*edit_duty.Out, error) {
	if in.Date == nil {
		return nil, errors.New("date is required")
	}

	// Get existing duty by date to obtain its ID
	existingDuty, err := u.r.GetDutyByDate(ctx, in.Date.String())
	if err != nil {
		return nil, err
	}

	arg := dbstore.UpdateDutyParams{
		ID:         existingDuty.ID,
		Date:       in.Date.String(),
		InWorkdays: in.InWorkdays,
		InHolidays: in.InHolidays,
	}

	if err := u.r.UpdateDuty(ctx, arg); err != nil {
		return nil, err
	}

	return &edit_duty.Out{Ok: true}, nil
}
