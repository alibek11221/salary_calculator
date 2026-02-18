package list

import (
	"context"

	"salary_calculator/internal/dto/list_duties"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context) (*list_duties.Out, error) {
	duties, err := u.r.ListDuties(ctx)
	if err != nil {
		return nil, err
	}

	outDuties := make([]list_duties.DutyItem, len(duties))

	for i, d := range duties {
		outDuties[i] = list_duties.DutyItem{
			Date:       d.Date,
			InWorkdays: d.InWorkdays,
			InHolidays: d.InHolidays,
		}
	}

	return &list_duties.Out{Duties: outDuties}, nil
}
