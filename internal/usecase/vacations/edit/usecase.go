package edit

import (
	"context"
	"errors"

	"salary_calculator/internal/dto/edit_vacation"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/database"
	"salary_calculator/internal/services/vacation_pay"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type usecase struct {
	r repo
}

func New(r repo) *usecase {
	return &usecase{r: r}
}

func (u *usecase) Do(ctx context.Context, in edit_vacation.In) (*edit_vacation.Out, error) {
	if err := vacation_pay.ValidatePeriod(in.From.Time, in.To.Time); err != nil {
		return nil, err
	}

	var id pgtype.UUID
	if err := id.Scan(in.ID); err != nil {
		return nil, ErrInvalidID
	}

	err := u.r.UpdateVacation(ctx, dbstore.UpdateVacationParams{
		ID:       id,
		DateFrom: in.From.PG(),
		DateTo:   in.To.PG(),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ExclusionViolationCode {
			return nil, ErrOverlappingVacation
		}

		return nil, err
	}

	return &edit_vacation.Out{Ok: true}, nil
}
