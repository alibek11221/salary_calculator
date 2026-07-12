package list_test

import (
	"context"
	"errors"
	"testing"

	list_bonuses_dto "salary_calculator/internal/dto/list_bonuses"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/usecase/bonuses/list"

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
	firstID := mustUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	secondID := mustUUID("6ba7b811-9dad-11d1-80b4-00c04fd430c8")

	tests := []struct {
		name       string
		setup      func(r *Mockrepo)
		want       *list_bonuses_dto.Out
		wantErrMsg string
	}{
		{
			name: "success maps all fields",
			setup: func(r *Mockrepo) {
				r.EXPECT().ListBonuses(gomock.Any()).Return([]dbstore.Bonuse{
					{ID: firstID, Value: 50000, Date: "2026_07"},
					{ID: secondID, Value: 12500.5, Date: "2026_08"},
				}, nil)
			},
			want: &list_bonuses_dto.Out{Bonuses: []list_bonuses_dto.Bonus{
				{ID: firstID.String(), Value: 50000, Date: mustDate("2026_07")},
				{ID: secondID.String(), Value: 12500.5, Date: mustDate("2026_08")},
			}},
		},
		{
			name: "empty list",
			setup: func(r *Mockrepo) {
				r.EXPECT().ListBonuses(gomock.Any()).Return([]dbstore.Bonuse{}, nil)
			},
			want: &list_bonuses_dto.Out{Bonuses: []list_bonuses_dto.Bonus{}},
		},
		{
			name: "invalid date in row",
			setup: func(r *Mockrepo) {
				r.EXPECT().ListBonuses(gomock.Any()).Return([]dbstore.Bonuse{
					{ID: firstID, Value: 50000, Date: "bad"},
				}, nil)
			},
			wantErrMsg: "invalid date format: expected 'year_month', got 'bad'",
		},
		{
			name: "db error passes through",
			setup: func(r *Mockrepo) {
				r.EXPECT().ListBonuses(gomock.Any()).Return(nil, errors.New("db error"))
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

			got, err := list.New(r).Do(context.Background())

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
