package get_bonuses_test

import (
	"context"
	"errors"
	"testing"

	"salary_calculator/internal/dto/get_bonuses"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	get_bonuses_uc "salary_calculator/internal/usecase/get_bonuses"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_Do(t *testing.T) {
	type fields struct {
		r *Mockrepo
	}

	sd, _ := value_objects.NewSalaryDate("2025_01")
	validID := "550e8400-e29b-41d4-a716-446655440000"
	var pgID pgtype.UUID
	pgID.Scan(validID)

	tests := []struct {
		name    string
		setup   func(f fields)
		want    *get_bonuses.Out
		wantErr bool
	}{
		{
			name: "success",
			setup: func(f fields) {
				f.r.EXPECT().ListBonuses(gomock.Any()).Return([]dbstore.Bonuse{
					{
						ID:          pgID,
						Value:       50000,
						Date:        "2025_01",
						Coefficient: 1.5,
					},
				}, nil)
			},
			want: &get_bonuses.Out{
				Bonuses: []get_bonuses.Bonus{
					{
						ID:          validID,
						Value:       50000,
						Date:        *sd,
						Coefficient: 1.5,
					},
				},
			},
		},
		{
			name: "db error",
			setup: func(f fields) {
				f.r.EXPECT().ListBonuses(gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "invalid date in db",
			setup: func(f fields) {
				f.r.EXPECT().ListBonuses(gomock.Any()).Return([]dbstore.Bonuse{
					{
						ID:          pgID,
						Value:       50000,
						Date:        "invalid_date",
						Coefficient: 1.5,
					},
				}, nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				r: NewMockrepo(ctrl),
			}
			if tt.setup != nil {
				tt.setup(f)
			}

			u := get_bonuses_uc.New(f.r)
			got, err := u.Do(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
