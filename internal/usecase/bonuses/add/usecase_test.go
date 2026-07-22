package add_test

import (
	"context"
	"errors"
	"testing"

	add_bonus_dto "salary_calculator/internal/dto/add_bonus"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/database"
	"salary_calculator/internal/usecase/bonuses/add"

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
	validIn := add_bonus_dto.In{Value: 50000, Date: mustDate("2026_07")}

	tests := []struct {
		name       string
		in         add_bonus_dto.In
		setup      func(r *Mockrepo)
		want       *add_bonus_dto.Out
		wantErrIs  error
		wantErrMsg string
	}{
		{
			name: "success",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertBonus(gomock.Any(), dbstore.InsertBonusParams{
					Value: 50000,
					Date:  "2026_07",
				}).Return(nil)
			},
			want: &add_bonus_dto.Out{Ok: true},
		},
		{
			name:       "nil date",
			in:         add_bonus_dto.In{Value: 50000},
			wantErrMsg: "date is required",
		},
		{
			name: "duplicate maps to domain error",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertBonus(gomock.Any(), gomock.Any()).
					Return(&pgconn.PgError{Code: database.DuplicateEntryCode})
			},
			wantErrIs: add.ErrDuplicateBonus,
		},
		{
			name: "other db error passes through",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertBonus(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
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
