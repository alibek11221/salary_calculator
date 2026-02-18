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

func TestUsecase_Do(t *testing.T) {
	type fields struct {
		r *Mockrepo
	}
	type args struct {
		ctx context.Context
		in  add_bonus_dto.In
	}

	sd, _ := value_objects.NewSalaryDate("2025_01")

	tests := []struct {
		name    string
		setup   func(f fields)
		args    args
		want    *add_bonus_dto.Out
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				in: add_bonus_dto.In{
					Value: 50000,
					Date:  sd,
				},
			},
			setup: func(f fields) {
				f.r.EXPECT().InsertBonus(gomock.Any(), dbstore.InsertBonusParams{
					Value:       50000,
					Date:        sd.String(),
					Coefficient: 1.5,
				}).Return(nil)
			},
			want: &add_bonus_dto.Out{Ok: true},
		},
		{
			name: "duplicate error",
			args: args{
				ctx: context.Background(),
				in: add_bonus_dto.In{
					Value: 50000,
					Date:  sd,
				},
			},
			setup: func(f fields) {
				f.r.EXPECT().InsertBonus(gomock.Any(), gomock.Any()).Return(&pgconn.PgError{Code: database.DuplicateEntryCode})
			},
			wantErr: add.ErrDuplicateBonus,
		},
		{
			name: "other error",
			args: args{
				ctx: context.Background(),
				in: add_bonus_dto.In{
					Value: 50000,
					Date:  sd,
				},
			},
			setup: func(f fields) {
				f.r.EXPECT().InsertBonus(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
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

			u := add.New(f.r)
			got, err := u.Do(tt.args.ctx, tt.args.in)

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
