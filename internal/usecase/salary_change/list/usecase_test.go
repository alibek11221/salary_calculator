package list_test

import (
	"context"
	"errors"
	"testing"

	"salary_calculator/internal/dto/list_salary_changes"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	get_salary_changes_uc "salary_calculator/internal/usecase/salary_change/list"

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
	_ = pgID.Scan(validID)

	tests := []struct {
		name    string
		setup   func(f fields)
		want    *list_salary_changes.Out
		wantErr bool
	}{
		{
			name: "success",
			setup: func(f fields) {
				f.r.EXPECT().ListChanges(gomock.Any()).Return([]dbstore.SalaryChange{
					{
						ID:         pgID,
						Salary:     100000,
						ChangeFrom: "2025_01",
					},
				}, nil)
			},
			want: &list_salary_changes.Out{
				Changes: []list_salary_changes.Change{
					{
						ID:    validID,
						Value: 100000,
						Date:  sd,
					},
				},
			},
		},
		{
			name: "db error",
			setup: func(f fields) {
				f.r.EXPECT().ListChanges(gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "invalid date in db",
			setup: func(f fields) {
				f.r.EXPECT().ListChanges(gomock.Any()).Return([]dbstore.SalaryChange{
					{
						ID:         pgID,
						Salary:     100000,
						ChangeFrom: "invalid_date",
					},
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "sorting by date",
			setup: func(f fields) {
				f.r.EXPECT().ListChanges(gomock.Any()).Return([]dbstore.SalaryChange{
					{
						ID:         pgID,
						Salary:     200000,
						ChangeFrom: "2025_02",
					},
					{
						ID:         pgID,
						Salary:     100000,
						ChangeFrom: "2025_01",
					},
					{
						ID:         pgID,
						Salary:     300000,
						ChangeFrom: "2024_12",
					},
				}, nil)
			},
			want: &list_salary_changes.Out{
				Changes: []list_salary_changes.Change{
					{
						ID:    validID,
						Value: 300000,
						Date:  func() *value_objects.SalaryDate { d, _ := value_objects.NewSalaryDate("2024_12"); return d }(),
					},
					{
						ID:    validID,
						Value: 100000,
						Date:  func() *value_objects.SalaryDate { d, _ := value_objects.NewSalaryDate("2025_01"); return d }(),
					},
					{
						ID:    validID,
						Value: 200000,
						Date:  func() *value_objects.SalaryDate { d, _ := value_objects.NewSalaryDate("2025_02"); return d }(),
					},
				},
			},
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

			u := get_salary_changes_uc.New(f.r)
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
