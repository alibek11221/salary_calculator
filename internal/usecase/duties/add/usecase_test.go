package add_test

import (
	"context"
	"errors"
	"testing"

	add_duty_dto "salary_calculator/internal/dto/add_duty"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/database"
	"salary_calculator/internal/usecase/duties/add"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func mustDate(s string) *value_objects.SalaryDate {
	d, err := value_objects.NewSalaryDate(s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestUsecase_Do(t *testing.T) {
	validIn := add_duty_dto.In{Date: mustDate("2026_07"), InWorkdays: 3, InHolidays: 2}

	tests := []struct {
		name       string
		in         add_duty_dto.In
		setup      func(r *Mockrepo)
		want       *add_duty_dto.Out
		wantErrIs  error  // sentinel-ошибка, проверяется через errors.Is
		wantErrMsg string // прочие ошибки, проверяется текст
	}{
		{
			name: "success",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertDuty(gomock.Any(), dbstore.InsertDutyParams{
					Date:       "2026_07",
					InWorkdays: 3,
					InHolidays: 2,
				}).Return(nil)
			},
			want: &add_duty_dto.Out{Ok: true},
		},
		{
			name:       "nil date",
			in:         add_duty_dto.In{InWorkdays: 3, InHolidays: 2},
			wantErrMsg: "date is required",
		},
		{
			name: "duplicate maps to domain error",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertDuty(gomock.Any(), gomock.Any()).
					Return(&pgconn.PgError{Code: database.DuplicateEntryCode})
			},
			wantErrIs: add.ErrDuplicateDuty,
		},
		{
			name: "other db error passes through",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertDuty(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErrMsg: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockrepo(ctrl)
			if tt.setup != nil {
				tt.setup(r)
			}

			got, err := add.New(r).Do(context.Background(), tt.in)

			switch {
			case tt.wantErrIs != nil:
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, got)
			case tt.wantErrMsg != "":
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, got)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
