package edit_test

import (
	"context"
	"errors"
	"testing"

	edit_bonus_dto "salary_calculator/internal/dto/edit_bonus"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/usecase/bonuses/edit"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

const bonusID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

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
	validIn := edit_bonus_dto.In{ID: bonusID, Value: 75000, Date: mustDate("2026_07")}

	tests := []struct {
		name       string
		in         edit_bonus_dto.In
		setup      func(r *Mockrepo)
		want       *edit_bonus_dto.Out
		wantErrMsg string
		wantAnyErr bool // текст не проверяем — зависит от внутренностей pgx
	}{
		{
			name: "success",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().UpdateBonus(gomock.Any(), dbstore.UpdateBonusParams{
					ID:    mustUUID(bonusID),
					Value: 75000,
					Date:  "2026_07",
				}).Return(nil)
			},
			want: &edit_bonus_dto.Out{Ok: true},
		},
		{
			name:       "nil date",
			in:         edit_bonus_dto.In{ID: bonusID, Value: 75000},
			wantErrMsg: "date is required",
		},
		{
			name:       "invalid id",
			in:         edit_bonus_dto.In{ID: "not-a-uuid", Value: 75000, Date: mustDate("2026_07")},
			wantAnyErr: true,
		},
		{
			name: "db error passes through",
			in:   validIn,
			setup: func(r *Mockrepo) {
				r.EXPECT().UpdateBonus(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
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

			got, err := edit.New(r).Do(context.Background(), tt.in)

			switch {
			case tt.wantAnyErr:
				assert.Error(t, err)
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
