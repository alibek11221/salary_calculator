package delete_test

import (
	"context"
	"errors"
	"testing"

	delete_duty_dto "salary_calculator/internal/dto/delete_duty"
	delete_duty_uc "salary_calculator/internal/usecase/duties/delete"

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
		in  delete_duty_dto.In
	}

	validID := "550e8400-e29b-41d4-a716-446655440000"
	var pgID pgtype.UUID
	pgID.Scan(validID)

	tests := []struct {
		name    string
		setup   func(f fields)
		args    args
		want    *delete_duty_dto.Out
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				in: delete_duty_dto.In{
					ID: validID,
				},
			},
			setup: func(f fields) {
				f.r.EXPECT().DeleteDuty(gomock.Any(), pgID).Return(nil)
			},
			want: &delete_duty_dto.Out{Ok: true},
		},
		{
			name: "invalid id",
			args: args{
				ctx: context.Background(),
				in: delete_duty_dto.In{
					ID: "invalid-uuid",
				},
			},
			wantErr: true,
		},
		{
			name: "db error",
			args: args{
				ctx: context.Background(),
				in: delete_duty_dto.In{
					ID: validID,
				},
			},
			setup: func(f fields) {
				f.r.EXPECT().DeleteDuty(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
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

			u := delete_duty_uc.New(f.r)
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
