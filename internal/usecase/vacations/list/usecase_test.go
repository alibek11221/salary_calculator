package list_test

import (
	"context"
	"testing"
	"time"

	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/usecase/vacations/list"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsecase_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := NewMockrepo(ctrl)
	r.EXPECT().ListVacations(gomock.Any()).Return([]dbstore.Vacation{
		{
			ID:       pgtype.UUID{Bytes: [16]byte{0x01, 0x98, 0xc5, 0xb6, 0x11, 0x11, 0x70, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, Valid: true},
			DateFrom: pgtype.Date{Time: time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC), Valid: true},
			DateTo:   pgtype.Date{Time: time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC), Valid: true},
		},
	}, nil)

	got, err := list.New(r).Do(context.Background())
	require.NoError(t, err)
	require.Len(t, got.Vacations, 1)
	assert.Equal(t, "0198c5b6-1111-7000-8000-000000000001", got.Vacations[0].ID)
	assert.Equal(t, 3, got.Vacations[0].From.Day())
	assert.Equal(t, 16, got.Vacations[0].To.Day())
}
