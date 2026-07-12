package list

import (
	"context"
	"fmt"

	"salary_calculator/internal/dto/list_vacations"
	"salary_calculator/internal/pkg/types"

	"github.com/jackc/pgx/v5/pgtype"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context) (*list_vacations.Out, error) {
	vacations, err := u.r.ListVacations(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]list_vacations.VacationItem, len(vacations))
	for i, v := range vacations {
		items[i] = list_vacations.VacationItem{
			ID:   uuidString(v.ID),
			From: types.ISODate{Time: v.DateFrom.Time},
			To:   types.ISODate{Time: v.DateTo.Time},
		}
	}

	return &list_vacations.Out{Vacations: items}, nil
}

func uuidString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
