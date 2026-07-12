package add_test

import (
	"context"
	"errors"
	"testing"

	add_vacation_dto "salary_calculator/internal/dto/add_vacation"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/database"
	"salary_calculator/internal/pkg/types"
	"salary_calculator/internal/services/vacation_pay"
	"salary_calculator/internal/usecase/vacations/add"

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
	validIn := add_vacation_dto.In{From: iso("2026-08-03"), To: iso("2026-08-16")}

	tests := []struct {
		name    string
		in      add_vacation_dto.In
		setup   func(r *Mockrepo)
		want    *add_vacation_dto.Out
		wantErr error
	}{
		{
			name: "success",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertVacation(gomock.Any(), dbstore.InsertVacationParams{
					DateFrom: iso("2026-08-03").PG(),
					DateTo:   iso("2026-08-16").PG(),
				}).Return(nil)
			},
			want: &add_vacation_dto.Out{Ok: true},
		},
		{
			name: "overlap maps to domain error",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertVacation(gomock.Any(), gomock.Any()).
					Return(&pgconn.PgError{Code: database.ExclusionViolationCode})
			},
			wantErr: add.ErrOverlappingVacation,
		},
		{
			name:    "from after to",
			in:      add_vacation_dto.In{From: iso("2026-08-16"), To: iso("2026-08-03")},
			wantErr: vacation_pay.ErrInvalidPeriod,
		},
		{
			name:    "year before start year",
			in:      add_vacation_dto.In{From: iso("2020-01-01"), To: iso("2020-01-05")},
			wantErr: vacation_pay.ErrInvalidPeriod,
		},
		{
			name: "other db error passes through",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().InsertVacation(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
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

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
