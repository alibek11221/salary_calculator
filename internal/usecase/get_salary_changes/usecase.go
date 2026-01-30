package get_salary_changes

import (
	"context"
	"salary_calculator/internal/dto/get_salary_changes"
	"salary_calculator/internal/dto/value_objects"
	"slices"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context) (*get_salary_changes.Out, error) {
	changes, err := u.r.ListChanges(ctx)
	if err != nil {
		return nil, err
	}

	outChanges := make([]get_salary_changes.Change, len(changes))

	for i, change := range changes {
		parsedDate, err := value_objects.NewSalaryDate(change.ChangeFrom)
		if err != nil {
			return nil, err
		}

		outChanges[i] = get_salary_changes.Change{
			ID:    change.ID.String(),
			Value: change.Salary,
			Date:  parsedDate,
		}
	}

	slices.SortFunc(outChanges, func(a, b get_salary_changes.Change) int {
		return a.Date.Compare(b.Date)
	})

	return &get_salary_changes.Out{Changes: outChanges}, nil
}
