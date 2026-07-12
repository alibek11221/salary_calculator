package estimate_test

import (
	"context"
	"testing"

	estimate_dto "salary_calculator/internal/dto/estimate_vacation"
	report_dto "salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/types"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/vacation_pay"
	"salary_calculator/internal/usecase/vacations/estimate"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func iso(s string) types.ISODate {
	d, err := types.ParseISODate(s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestUsecase_Do(t *testing.T) {
	t.Run("vacation crossing month boundary", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		r := NewMockrepo(ctrl)
		vp := NewMockvacationPayService(ctrl)
		rep := NewMockreportUsecase(ctrl)

		in := estimate_dto.In{From: iso("2026-08-24"), To: iso("2026-09-06")}

		earnings := &vacation_pay.EarningsData{}
		vp.EXPECT().LoadEarningsData(gomock.Any()).Times(1).Return(earnings, nil)
		vp.EXPECT().CalculatePayWith(earnings, in.From.Time, in.To.Time).
			Return(&vacation_pay.Pay{AvgDaily: 10000, PaidDays: 14, Gross: 140000}, nil)
		r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), "2026_08").
			Return(dbstore.SalaryChange{Salary: 100000}, nil) // НДФЛ 13%

		// Порядок: август, потом сентябрь. В каждый отчёт прокидываются
		// уже загруженные earnings — повторных выборок нет.
		gomock.InOrder(
			rep.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, rin report_dto.In) (*report_dto.Out, error) {
				assert.Equal(t, 2026, rin.Year)
				assert.Equal(t, 8, rin.Month)
				require.Len(t, rin.ExtraVacations, 1)
				assert.Equal(t, in.From.Time, rin.ExtraVacations[0].DateFrom.Time)
				assert.Same(t, earnings, rin.Earnings)
				return &report_dto.Out{Result: &calculator.SalaryCalculationResult{Advance: 1, Salary: 2, Total: 100}}, nil
			}),
			rep.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, rin report_dto.In) (*report_dto.Out, error) {
				assert.Equal(t, 2026, rin.Year)
				assert.Equal(t, 9, rin.Month)
				assert.Same(t, earnings, rin.Earnings)
				return &report_dto.Out{Result: &calculator.SalaryCalculationResult{Advance: 3, Salary: 4, Total: 200}}, nil
			}),
		)

		got, err := estimate.New(r, vp, rep).Do(context.Background(), in)

		require.NoError(t, err)
		assert.InDelta(t, 10000.0, got.AvgDaily, 0.01)
		assert.Equal(t, 14, got.PaidDays)
		assert.InDelta(t, 140000.0, got.Gross, 0.01)
		assert.InDelta(t, 121800.0, got.Net, 0.01) // 140000 − 13%
		require.Len(t, got.AffectedMonths, 2)
		assert.Equal(t, "2026_08", got.AffectedMonths[0].Month)
		assert.Equal(t, 100.0, got.AffectedMonths[0].Total)
		assert.Equal(t, "2026_09", got.AffectedMonths[1].Month)
		assert.Equal(t, 200.0, got.AffectedMonths[1].Total)
	})

	t.Run("invalid period", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, err := estimate.New(NewMockrepo(ctrl), NewMockvacationPayService(ctrl), NewMockreportUsecase(ctrl)).
			Do(context.Background(), estimate_dto.In{From: iso("2026-08-16"), To: iso("2026-08-03")})
		assert.ErrorIs(t, err, vacation_pay.ErrInvalidPeriod)
	})

	t.Run("earnings load error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		vp := NewMockvacationPayService(ctrl)
		vp.EXPECT().LoadEarningsData(gomock.Any()).Return(nil, assert.AnError)

		_, err := estimate.New(NewMockrepo(ctrl), vp, NewMockreportUsecase(ctrl)).
			Do(context.Background(), estimate_dto.In{From: iso("2026-08-03"), To: iso("2026-08-16")})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("pay calc error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		vp := NewMockvacationPayService(ctrl)
		vp.EXPECT().LoadEarningsData(gomock.Any()).Return(&vacation_pay.EarningsData{}, nil)
		vp.EXPECT().CalculatePayWith(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, vacation_pay.ErrNoEarningsData)

		_, err := estimate.New(NewMockrepo(ctrl), vp, NewMockreportUsecase(ctrl)).
			Do(context.Background(), estimate_dto.In{From: iso("2026-08-03"), To: iso("2026-08-16")})
		assert.ErrorIs(t, err, vacation_pay.ErrNoEarningsData)
	})
}
