package list

import (
	"context"
	"slices"

	"salary_calculator/internal/dto/list_salary_changes"
	"salary_calculator/internal/dto/value_objects"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context) (*list_salary_changes.Out, error) {
	changes, err := u.r.ListChanges(ctx)
	if err != nil {
		return nil, err
	}

	outChanges := make([]list_salary_changes.Change, len(changes))

	for i, change := range changes {
		parsedDate, err := value_objects.NewSalaryDate(change.ChangeFrom)
		if err != nil {
			return nil, err
		}

		outChanges[i] = list_salary_changes.Change{
			ID:    change.ID.String(),
			Value: change.Salary,
			Date:  parsedDate,
		}
	}

	slices.SortFunc(outChanges, func(a, b list_salary_changes.Change) int {
		if a.Date == nil && b.Date == nil {
			return 0
		}
		if a.Date == nil {
			return -1
		}
		return a.Date.Compare(b.Date)
	})

	return &list_salary_changes.Out{Changes: outChanges}, nil
}
