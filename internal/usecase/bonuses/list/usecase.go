package list

import (
	"context"

	"salary_calculator/internal/dto/list_bonuses"
	"salary_calculator/internal/dto/value_objects"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context) (*list_bonuses.Out, error) {
	bonuses, err := u.r.ListBonuses(ctx)
	if err != nil {
		return nil, err
	}

	outBonuses := make([]list_bonuses.Bonus, len(bonuses))

	for i, b := range bonuses {
		parsedDate, err := value_objects.NewSalaryDate(b.Date)
		if err != nil {
			return nil, err
		}

		outBonuses[i] = list_bonuses.Bonus{
			ID:    b.ID.String(),
			Value: b.Value,
			Date:  parsedDate,
		}
	}

	return &list_bonuses.Out{Bonuses: outBonuses}, nil
}
