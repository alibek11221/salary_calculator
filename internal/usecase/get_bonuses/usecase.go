package get_bonuses

import (
	"context"

	"salary_calculator/internal/dto/get_bonuses"
	"salary_calculator/internal/dto/value_objects"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context) (*get_bonuses.Out, error) {
	bonuses, err := u.r.ListBonuses(ctx)
	if err != nil {
		return nil, err
	}

	outBonuses := make([]get_bonuses.Bonus, len(bonuses))

	for i, b := range bonuses {
		parsedDate, err := value_objects.NewSalaryDate(b.Date)
		if err != nil {
			return nil, err
		}

		outBonuses[i] = get_bonuses.Bonus{
			ID:          b.ID.String(),
			Value:       b.Value,
			Date:        *parsedDate,
			Coefficient: b.Coefficient,
		}
	}

	return &get_bonuses.Out{Bonuses: outBonuses}, nil
}
