package edit_salary_change_test

import (
	"context"
	"errors"
	"testing"

	edit_salary_change_dto "salary_calculator/internal/dto/edit_salary_change"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	edit_salary_change_uc "salary_calculator/internal/usecase/edit_salary_change"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_Do(t *testing.T) {
	type fields struct {
		r *Mockrepo
	}
	type args struct {
		ctx context.Context
		in  edit_salary_change_dto.In
	}

	sd, _ := value_objects.NewSalaryDate("2025_01")
	validID := "550e8400-e29b-41d4-a716-446655440000"
	var pgID pgtype.UUID
	_ = pgID.Scan(validID)

	tests := []struct {
		name    string
		setup   func(f fields)
		args    args
		want    *edit_salary_change_dto.Out
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				in: edit_salary_change_dto.In{
					ID:    validID,
					Value: 120000,
					Date:  sd,
				},
			},
			setup: func(f fields) {
				f.r.EXPECT().UpdateChange(gomock.Any(), dbstore.UpdateChangeParams{
					ID:         pgID,
					Salary:     120000,
					ChangeFrom: sd.String(),
				}).Return(nil)
			},
			want: &edit_salary_change_dto.Out{Ok: true},
		},
		{
			name: "invalid id",
			args: args{
				ctx: context.Background(),
				in: edit_salary_change_dto.In{
					ID: "invalid-uuid",
				},
			},
			wantErr: true,
		},
		{
			name: "db error",
			args: args{
				ctx: context.Background(),
				in: edit_salary_change_dto.In{
					ID:    validID,
					Value: 120000,
					Date:  sd,
				},
			},
			setup: func(f fields) {
				f.r.EXPECT().UpdateChange(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
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

			u := edit_salary_change_uc.New(f.r)
			got, err := u.Do(tt.args.ctx, tt.args.in)

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
