package delete_test

import (
	"context"
	"testing"

	delete_vacation_dto "salary_calculator/internal/dto/delete_vacation"
	del "salary_calculator/internal/usecase/vacations/delete"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_Do(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockrepo(ctrl)
		r.EXPECT().DeleteVacation(gomock.Any(), gomock.Any()).Return(nil)

		got, err := del.New(r).Do(context.Background(), delete_vacation_dto.In{
			ID: "0198c5b6-1111-7000-8000-000000000001",
		})
		assert.NoError(t, err)
		assert.True(t, got.Ok)
	})

	t.Run("bad id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, err := del.New(NewMockrepo(ctrl)).Do(context.Background(), delete_vacation_dto.In{ID: "nope"})
		assert.ErrorIs(t, err, del.ErrInvalidID)
	})
}
