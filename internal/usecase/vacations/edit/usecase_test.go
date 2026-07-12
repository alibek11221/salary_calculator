package edit_test

import (
	"context"
	"testing"

	edit_vacation_dto "salary_calculator/internal/dto/edit_vacation"
	"salary_calculator/internal/pkg/database"
	"salary_calculator/internal/pkg/types"
	"salary_calculator/internal/usecase/vacations/edit"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func iso(s string) types.ISODate {
	d, err := types.ParseISODate(s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestUsecase_Do(t *testing.T) {
	const validID = "0198c5b6-1111-7000-8000-000000000001"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockrepo(ctrl)
		r.EXPECT().UpdateVacation(gomock.Any(), gomock.Any()).Return(nil)

		got, err := edit.New(r).Do(context.Background(), edit_vacation_dto.In{
			ID: validID, From: iso("2026-08-03"), To: iso("2026-08-16"),
		})
		assert.NoError(t, err)
		assert.True(t, got.Ok)
	})

	t.Run("bad id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, err := edit.New(NewMockrepo(ctrl)).Do(context.Background(), edit_vacation_dto.In{
			ID: "not-a-uuid", From: iso("2026-08-03"), To: iso("2026-08-16"),
		})
		assert.ErrorIs(t, err, edit.ErrInvalidID)
	})

	t.Run("overlap", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockrepo(ctrl)
		r.EXPECT().UpdateVacation(gomock.Any(), gomock.Any()).
			Return(&pgconn.PgError{Code: database.ExclusionViolationCode})

		_, err := edit.New(r).Do(context.Background(), edit_vacation_dto.In{
			ID: validID, From: iso("2026-08-03"), To: iso("2026-08-16"),
		})
		assert.ErrorIs(t, err, edit.ErrOverlappingVacation)
	})
}
