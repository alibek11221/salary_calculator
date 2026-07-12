package delete_test

import (
	"context"
	"errors"
	"testing"

	delete_duty_dto "salary_calculator/internal/dto/delete_duty"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	delete_duty_uc "salary_calculator/internal/usecase/duties/delete"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func mustDate(s string) *value_objects.SalaryDate {
	d, err := value_objects.NewSalaryDate(s)
	if err != nil {
		panic(err)
	}
	return d
}

func mustUUID(s string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(s); err != nil {
		panic(err)
	}
	return id
}

func TestUsecase_Do(t *testing.T) {
	dutyID := mustUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	validIn := delete_duty_dto.In{Date: mustDate("2026_07")}

	tests := []struct {
		name       string
		in         delete_duty_dto.In
		setup      func(r *Mockrepo)
		want       *delete_duty_dto.Out
		wantErrMsg string
	}{
		{
			name: "success",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().GetDutyByDate(gomock.Any(), "2026_07").
					Return(dbstore.Duty{ID: dutyID, Date: "2026_07"}, nil)
				r.EXPECT().DeleteDuty(gomock.Any(), dutyID).Return(nil)
			},
			want: &delete_duty_dto.Out{Ok: true},
		},
		{
			name:       "nil date",
			in:         delete_duty_dto.In{},
			wantErrMsg: "date is required",
		},
		{
			name: "duty lookup error",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().GetDutyByDate(gomock.Any(), "2026_07").
					Return(dbstore.Duty{}, errors.New("not found"))
			},
			wantErrMsg: "not found",
		},
		{
			name: "delete db error passes through",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().GetDutyByDate(gomock.Any(), "2026_07").
					Return(dbstore.Duty{ID: dutyID, Date: "2026_07"}, nil)
				r.EXPECT().DeleteDuty(gomock.Any(), dutyID).Return(errors.New("db error"))
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

			got, err := delete_duty_uc.New(r).Do(context.Background(), tt.in)

			if tt.wantErrMsg != "" {
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
