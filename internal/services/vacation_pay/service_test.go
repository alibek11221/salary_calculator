package vacation_pay_test

import (
	"context"
	"testing"
	"time"

	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/http/work_calendar_parser"
	"salary_calculator/internal/pkg/types"
	"salary_calculator/internal/services/vacation_pay"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func day(y int, m time.Month, d, typeID int) work_calendar_parser.Day {
	return work_calendar_parser.Day{
		Date:   types.Date{Time: time.Date(y, m, d, 0, 0, 0, 0, time.UTC)},
		TypeId: typeID,
	}
}

func pgDate(y int, m time.Month, d int) pgtype.Date {
	return pgtype.Date{Time: time.Date(y, m, d, 0, 0, 0, 0, time.UTC), Valid: true}
}

func vacation(fy int, fm time.Month, fd, ty int, tm time.Month, td int) dbstore.Vacation {
	return dbstore.Vacation{
		ID:       pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		DateFrom: pgDate(fy, fm, fd),
		DateTo:   pgDate(ty, tm, td),
	}
}

// Июнь 2026 (реалистичный фрагмент): 1–5 рабочие (пн–пт), 6–7 выходные,
// 8–11 рабочие, 12 — госпраздник (День России), 13–14 выходные, 15–19 рабочие.
func juneDays() []work_calendar_parser.Day {
	days := []work_calendar_parser.Day{}
	workdays := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 8: true, 9: true, 10: true, 11: true, 15: true, 16: true, 17: true, 18: true, 19: true}
	for d := 1; d <= 19; d++ {
		typeID := 2 // выходной
		if workdays[d] {
			typeID = 1
		}
		if d == 12 {
			typeID = 3 // праздник
		}
		days = append(days, day(2026, time.June, d, typeID))
	}
	return days
}

func TestPaidDays(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parser := NewMockcalendarParser(ctrl)
	parser.EXPECT().Parse(2026, 6).AnyTimes().Return(&work_calendar_parser.WorkdayResponse{Days: juneDays()}, nil)

	svc := vacation_pay.New(NewMockrepo(ctrl), parser)

	t.Run("holiday excluded, weekends included", func(t *testing.T) {
		// 8–14 июня: 7 календарных дней, из них 12-е — праздник → 6 оплачиваемых.
		got, err := svc.PaidDays(
			time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 6, 14, 0, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)
		assert.Equal(t, 6, got)
	})

	t.Run("no holidays in range", func(t *testing.T) {
		// 1–7 июня: 7 календарных, праздников нет → 7.
		got, err := svc.PaidDays(
			time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 6, 7, 0, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)
		assert.Equal(t, 7, got)
	})

	t.Run("to before from is error", func(t *testing.T) {
		_, err := svc.PaidDays(
			time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC),
		)
		assert.ErrorIs(t, err, vacation_pay.ErrInvalidPeriod)
	})

	t.Run("missing calendar propagates error", func(t *testing.T) {
		parser.EXPECT().Parse(2031, 1).Return(nil, assert.AnError)
		_, err := svc.PaidDays(
			time.Date(2031, 1, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2031, 1, 6, 0, 0, 0, 0, time.UTC),
		)
		assert.Error(t, err)
	})

	t.Run("day missing from calendar is error", func(t *testing.T) {
		// juneDays покрывает только 1–19 июня → 20-е не найдено.
		_, err := svc.PaidDays(
			time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC),
		)
		assert.ErrorContains(t, err, "day 2026-06-20 not found in work calendar")
	})

	t.Run("period crossing month boundary parses each month once", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := NewMockcalendarParser(ctrl)
		p.EXPECT().Parse(2026, 6).Times(1).Return(&work_calendar_parser.WorkdayResponse{Days: []work_calendar_parser.Day{
			day(2026, time.June, 29, 1), day(2026, time.June, 30, 1),
		}}, nil)
		p.EXPECT().Parse(2026, 7).Times(1).Return(&work_calendar_parser.WorkdayResponse{Days: []work_calendar_parser.Day{
			day(2026, time.July, 1, 1), day(2026, time.July, 2, 3), // 2 июля — «праздник»
		}}, nil)

		got, err := vacation_pay.New(NewMockrepo(ctrl), p).PaidDays(
			time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, got) // 29, 30 июня и 1 июля; праздник 2 июля не оплачивается
	})
}

func TestWorkdayDeduction(t *testing.T) {
	tests := []struct {
		name      string
		vacations []dbstore.Vacation
		want      vacation_pay.Deduction
	}{
		{
			name:      "no vacations",
			vacations: nil,
			want:      vacation_pay.Deduction{},
		},
		{
			name:      "vacation in first half over weekend",
			vacations: []dbstore.Vacation{vacation(2026, time.June, 3, 2026, time.June, 9)},
			// Рабочие в периоде 3–9: 3,4,5,8,9 → 5, все до 15-го.
			want: vacation_pay.Deduction{FirstHalf: 5},
		},
		{
			name:      "vacation crossing month halves",
			vacations: []dbstore.Vacation{vacation(2026, time.June, 10, 2026, time.June, 17)},
			// Рабочие 10,11,15 — первая половина (12 праздник, 13–14 выходные); 16,17 — вторая.
			want: vacation_pay.Deduction{FirstHalf: 3, SecondHalf: 2},
		},
		{
			name: "two vacations sum up",
			vacations: []dbstore.Vacation{
				vacation(2026, time.June, 1, 2026, time.June, 2),
				vacation(2026, time.June, 18, 2026, time.June, 19),
			},
			want: vacation_pay.Deduction{FirstHalf: 2, SecondHalf: 2},
		},
		{
			name: "vacation from another month ignored",
			// Отпуск целиком в мае — июньских дней не трогает.
			vacations: []dbstore.Vacation{vacation(2026, time.May, 1, 2026, time.May, 10)},
			want:      vacation_pay.Deduction{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vacation_pay.WorkdayDeduction(juneDays(), tt.vacations)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLoadEarningsData(t *testing.T) {
	t.Run("loads changes and bonuses once", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		changes := []dbstore.SalaryChange{{Salary: 100000, ChangeFrom: "2024_01"}}
		bonuses := []dbstore.Bonuse{{Value: 50000, Date: "2025_03"}}

		r := NewMockrepo(ctrl)
		r.EXPECT().ListChanges(gomock.Any()).Times(1).Return(changes, nil)
		r.EXPECT().ListBonuses(gomock.Any()).Times(1).Return(bonuses, nil)

		data, err := vacation_pay.New(r, NewMockcalendarParser(ctrl)).LoadEarningsData(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, changes, data.Changes)
		assert.Equal(t, bonuses, data.Bonuses)
	})

	t.Run("changes error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		r := NewMockrepo(ctrl)
		r.EXPECT().ListChanges(gomock.Any()).Return(nil, assert.AnError)

		_, err := vacation_pay.New(r, NewMockcalendarParser(ctrl)).LoadEarningsData(context.Background())
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("bonuses error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		r := NewMockrepo(ctrl)
		r.EXPECT().ListChanges(gomock.Any()).Return(nil, nil)
		r.EXPECT().ListBonuses(gomock.Any()).Return(nil, assert.AnError)

		_, err := vacation_pay.New(r, NewMockcalendarParser(ctrl)).LoadEarningsData(context.Background())
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestAverageDailyEarnings(t *testing.T) {
	start := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC) // расчётный период 2024_08..2025_07

	tests := []struct {
		name    string
		data    *vacation_pay.EarningsData
		want    float64
		wantErr error
	}{
		{
			name: "full 12 months with raise and bonus",
			data: &vacation_pay.EarningsData{
				Changes: []dbstore.SalaryChange{
					{Salary: 100000, ChangeFrom: "2024_01"},
					{Salary: 200000, ChangeFrom: "2025_06"},
				},
				Bonuses: []dbstore.Bonuse{
					{Value: 50000, Date: "2025_03"},
					{Value: 99999, Date: "2023_01"}, // вне периода — игнор
				},
			},
			// 10 мес × 100000 + 2 мес × 200000 + 50000 = 1450000
			want: 1450000.0 / 12.0 / 29.3,
		},
		{
			name: "short history divides by actual months",
			data: &vacation_pay.EarningsData{
				Changes: []dbstore.SalaryChange{
					{Salary: 200000, ChangeFrom: "2025_06"},
				},
			},
			// Данные только за 2025_06 и 2025_07 → 400000 / 2 / 29.3
			want: 400000.0 / 2.0 / 29.3,
		},
		{
			name: "duplicate change_from keeps first occurrence",
			data: &vacation_pay.EarningsData{
				Changes: []dbstore.SalaryChange{
					{Salary: 100000, ChangeFrom: "2025_06"},
					{Salary: 999999, ChangeFrom: "2025_06"}, // дубликат — игнор
				},
			},
			// 2025_06 и 2025_07 по 100000 → 200000 / 2 / 29.3
			want: 200000.0 / 2.0 / 29.3,
		},
		{
			name:    "no data at all",
			data:    &vacation_pay.EarningsData{},
			wantErr: vacation_pay.ErrNoEarningsData,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: vacation_pay.ErrNoEarningsData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := vacation_pay.New(nil, nil)
			got, err := svc.AverageDailyEarnings(tt.data, start)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.InDelta(t, tt.want, got, 0.0001)
		})
	}
}

func TestCalculatePayWith(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parser := NewMockcalendarParser(ctrl)
	parser.EXPECT().Parse(2026, 6).AnyTimes().Return(&work_calendar_parser.WorkdayResponse{Days: juneDays()}, nil)

	svc := vacation_pay.New(NewMockrepo(ctrl), parser) // repo не должен трогаться

	data := &vacation_pay.EarningsData{
		Changes: []dbstore.SalaryChange{
			{Salary: 293000, ChangeFrom: "2024_01"}, // 293000/29.3 = 10000/день ровно
		},
	}

	// 8–14 июня 2026: 6 оплачиваемых дней (12-е — праздник).
	pay, err := svc.CalculatePayWith(data,
		time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 14, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)
	assert.InDelta(t, 10000.0, pay.AvgDaily, 0.0001)
	assert.Equal(t, 6, pay.PaidDays)
	assert.InDelta(t, 60000.0, pay.Gross, 0.001)
}

func TestCalculatePay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := NewMockrepo(ctrl)
	r.EXPECT().ListChanges(gomock.Any()).Return([]dbstore.SalaryChange{
		{Salary: 293000, ChangeFrom: "2024_01"}, // 293000/29.3 = 10000/день ровно
	}, nil)
	r.EXPECT().ListBonuses(gomock.Any()).Return(nil, nil)

	parser := NewMockcalendarParser(ctrl)
	parser.EXPECT().Parse(2026, 6).AnyTimes().Return(&work_calendar_parser.WorkdayResponse{Days: juneDays()}, nil)

	svc := vacation_pay.New(r, parser)

	// 8–14 июня 2026: 6 оплачиваемых дней (12-е — праздник).
	pay, err := svc.CalculatePay(context.Background(),
		time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 14, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)
	assert.InDelta(t, 10000.0, pay.AvgDaily, 0.0001)
	assert.Equal(t, 6, pay.PaidDays)
	assert.InDelta(t, 60000.0, pay.Gross, 0.001)
}
