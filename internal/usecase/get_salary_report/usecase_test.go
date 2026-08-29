package get_salary_report_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	report_dto "salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/http/work_calendar_parser"
	"salary_calculator/internal/pkg/types"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/vacation_pay"
	"salary_calculator/internal/services/work_days"
	uc "salary_calculator/internal/usecase/get_salary_report"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// Июнь 2026, дни 1–8: 1–5 рабочие, 6–7 выходные, 8 рабочий.
func juneCalendar() *work_calendar_parser.WorkdayResponse {
	return &work_calendar_parser.WorkdayResponse{
		Statistics: work_calendar_parser.Statistics{WorkDays: 20},
		Days: []work_calendar_parser.Day{
			day(2026, time.June, 1, 1), day(2026, time.June, 2, 1), day(2026, time.June, 3, 1),
			day(2026, time.June, 4, 1), day(2026, time.June, 5, 1), day(2026, time.June, 6, 2),
			day(2026, time.June, 7, 2), day(2026, time.June, 8, 1),
		},
	}
}

type fields struct {
	r       *Mockrepo
	parser  *MockworkdaysParser
	wdCalc  *MockworkdaysCalculator
	salCalc *MocksalaryCalculator
	vacPay  *MockvacationPay
}

func TestUsecase_Do(t *testing.T) {
	in := report_dto.In{Year: 2026, Month: 6}

	newFields := func(ctrl *gomock.Controller) fields {
		return fields{
			r:       NewMockrepo(ctrl),
			parser:  NewMockworkdaysParser(ctrl),
			wdCalc:  NewMockworkdaysCalculator(ctrl),
			salCalc: NewMocksalaryCalculator(ctrl),
			vacPay:  NewMockvacationPay(ctrl),
		}
	}

	t.Run("success without vacations", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFields(ctrl)

		f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), "2026_06").
			Return(dbstore.SalaryChange{Salary: 100000, ChangeFrom: "2026_01"}, nil)
		f.r.EXPECT().GetVacationsInRange(gomock.Any(), gomock.Any()).Return(nil, nil)
		f.parser.EXPECT().Parse(2026, 6).Return(juneCalendar(), nil)
		f.wdCalc.EXPECT().CalculateWorkDaysForMonth(gomock.Any()).
			Return(&work_days.WorkdaysForMonth{TotalWorkdays: 20, FirstHalfDays: 10, SecondHalfDays: 10})
		f.vacPay.EXPECT().NetPaymentsForMonth(gomock.Any(), []dbstore.Vacation{}, 2026, 6, 13.0, nil).
			Return(nil, nil)

		var gotCtx value_objects.SalaryCalculationContext
		f.salCalc.EXPECT().CalculateSalary(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ value_objects.SalaryDate, sCtx value_objects.SalaryCalculationContext) (*calculator.SalaryCalculationResult, error) {
				gotCtx = sCtx
				return &calculator.SalaryCalculationResult{InHand: calculator.PayoutBreakdown{Total: 1}}, nil
			})

		u := uc.New(f.r, f.parser, f.wdCalc, f.salCalc, f.vacPay)
		out, err := u.Do(context.Background(), in)

		require.NoError(t, err)
		assert.Equal(t, 100000.0, out.BaseSalary)
		assert.Equal(t, 10, gotCtx.Workdays().FirstHalfDays)
		assert.Empty(t, gotCtx.VacationPayments())
	})

	t.Run("success with vacation reduces days and adds payment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFields(ctrl)

		vac := dbstore.Vacation{
			ID:       pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
			DateFrom: pgDate(2026, time.June, 2),
			DateTo:   pgDate(2026, time.June, 8),
		}

		f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), "2026_06").
			Return(dbstore.SalaryChange{Salary: 100000, ChangeFrom: "2026_01"}, nil)
		f.r.EXPECT().GetVacationsInRange(gomock.Any(), dbstore.GetVacationsInRangeParams{
			RangeStart: pgDate(2026, time.June, 1),
			RangeEnd:   pgDate(2026, time.June, 30),
		}).Return([]dbstore.Vacation{vac}, nil)
		f.parser.EXPECT().Parse(2026, 6).Return(juneCalendar(), nil)
		f.wdCalc.EXPECT().CalculateWorkDaysForMonth(gomock.Any()).
			Return(&work_days.WorkdaysForMonth{TotalWorkdays: 20, FirstHalfDays: 10, SecondHalfDays: 10})
		payment := value_objects.ExtraPayment{
			Name:  vacation_pay.PaymentName,
			Value: 30450,
			T:     value_objects.Extra,
		}
		f.vacPay.EXPECT().NetPaymentsForMonth(gomock.Any(), []dbstore.Vacation{vac}, 2026, 6, 13.0, nil).
			Return([]value_objects.ExtraPayment{payment}, nil)

		var gotCtx value_objects.SalaryCalculationContext
		f.salCalc.EXPECT().CalculateSalary(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ value_objects.SalaryDate, sCtx value_objects.SalaryCalculationContext) (*calculator.SalaryCalculationResult, error) {
				gotCtx = sCtx
				return &calculator.SalaryCalculationResult{InHand: calculator.PayoutBreakdown{Total: 1}}, nil
			})

		u := uc.New(f.r, f.parser, f.wdCalc, f.salCalc, f.vacPay)
		_, err := u.Do(context.Background(), in)

		require.NoError(t, err)
		// Отпуск 2–8 июня покрывает рабочие 2,3,4,5,8 → все в первой половине.
		assert.Equal(t, 5, gotCtx.Workdays().FirstHalfDays)
		assert.Equal(t, 10, gotCtx.Workdays().SecondHalfDays)
		require.Len(t, gotCtx.VacationPayments(), 1)
		assert.Equal(t, payment, gotCtx.VacationPayments()[0])
	})

	t.Run("pre-supplied earnings forwarded to vacation pay", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFields(ctrl)

		vac := dbstore.Vacation{
			ID:       pgtype.UUID{Bytes: [16]byte{3}, Valid: true},
			DateFrom: pgDate(2026, time.June, 2),
			DateTo:   pgDate(2026, time.June, 3),
		}

		f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), "2026_06").
			Return(dbstore.SalaryChange{Salary: 100000, ChangeFrom: "2026_01"}, nil)
		f.r.EXPECT().GetVacationsInRange(gomock.Any(), gomock.Any()).Return([]dbstore.Vacation{vac}, nil)
		f.parser.EXPECT().Parse(2026, 6).Return(juneCalendar(), nil)
		f.wdCalc.EXPECT().CalculateWorkDaysForMonth(gomock.Any()).
			Return(&work_days.WorkdaysForMonth{TotalWorkdays: 20, FirstHalfDays: 10, SecondHalfDays: 10})

		earnings := &vacation_pay.EarningsData{}
		f.vacPay.EXPECT().NetPaymentsForMonth(gomock.Any(), []dbstore.Vacation{vac}, 2026, 6, 13.0, earnings).
			Return(nil, nil)

		f.salCalc.EXPECT().CalculateSalary(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&calculator.SalaryCalculationResult{InHand: calculator.PayoutBreakdown{Total: 1}}, nil)

		u := uc.New(f.r, f.parser, f.wdCalc, f.salCalc, f.vacPay)
		_, err := u.Do(context.Background(), report_dto.In{Year: 2026, Month: 6, Earnings: earnings})

		require.NoError(t, err)
	})

	t.Run("vacation tail from previous month still reduces days", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFields(ctrl)

		vac := dbstore.Vacation{
			ID:       pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
			DateFrom: pgDate(2026, time.May, 25),
			DateTo:   pgDate(2026, time.June, 3),
		}

		f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), gomock.Any()).
			Return(dbstore.SalaryChange{Salary: 100000, ChangeFrom: "2026_01"}, nil)
		f.r.EXPECT().GetVacationsInRange(gomock.Any(), gomock.Any()).Return([]dbstore.Vacation{vac}, nil)
		f.parser.EXPECT().Parse(2026, 6).Return(juneCalendar(), nil)
		f.wdCalc.EXPECT().CalculateWorkDaysForMonth(gomock.Any()).
			Return(&work_days.WorkdaysForMonth{TotalWorkdays: 20, FirstHalfDays: 10, SecondHalfDays: 10})
		f.vacPay.EXPECT().NetPaymentsForMonth(gomock.Any(), []dbstore.Vacation{vac}, 2026, 6, 13.0, nil).
			Return(nil, nil)

		var gotCtx value_objects.SalaryCalculationContext
		f.salCalc.EXPECT().CalculateSalary(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ value_objects.SalaryDate, sCtx value_objects.SalaryCalculationContext) (*calculator.SalaryCalculationResult, error) {
				gotCtx = sCtx
				return &calculator.SalaryCalculationResult{InHand: calculator.PayoutBreakdown{Total: 1}}, nil
			})

		u := uc.New(f.r, f.parser, f.wdCalc, f.salCalc, f.vacPay)
		_, err := u.Do(context.Background(), in)

		require.NoError(t, err)
		// Хвост отпуска съел рабочие 1,2,3 июня, но выплата была в мае.
		assert.Equal(t, 7, gotCtx.Workdays().FirstHalfDays)
		assert.Empty(t, gotCtx.VacationPayments())
	})

	t.Run("vacation payments error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFields(ctrl)

		f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), gomock.Any()).
			Return(dbstore.SalaryChange{Salary: 100000, ChangeFrom: "2026_01"}, nil)
		f.r.EXPECT().GetVacationsInRange(gomock.Any(), gomock.Any()).Return(nil, nil)
		f.parser.EXPECT().Parse(2026, 6).Return(juneCalendar(), nil)
		f.wdCalc.EXPECT().CalculateWorkDaysForMonth(gomock.Any()).
			Return(&work_days.WorkdaysForMonth{TotalWorkdays: 20, FirstHalfDays: 10, SecondHalfDays: 10})
		f.vacPay.EXPECT().NetPaymentsForMonth(gomock.Any(), gomock.Any(), 2026, 6, 13.0, nil).
			Return(nil, assert.AnError)

		u := uc.New(f.r, f.parser, f.wdCalc, f.salCalc, f.vacPay)
		_, err := u.Do(context.Background(), in)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("no salary change found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFields(ctrl)

		f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), gomock.Any()).
			Return(dbstore.SalaryChange{}, sql.ErrNoRows)
		f.r.EXPECT().GetVacationsInRange(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
		f.parser.EXPECT().Parse(gomock.Any(), gomock.Any()).AnyTimes().Return(juneCalendar(), nil)

		u := uc.New(f.r, f.parser, f.wdCalc, f.salCalc, f.vacPay)
		_, err := u.Do(context.Background(), in)
		assert.Error(t, err)
	})

	t.Run("vacations repo error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFields(ctrl)

		f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), gomock.Any()).AnyTimes().
			Return(dbstore.SalaryChange{Salary: 100000}, nil)
		f.r.EXPECT().GetVacationsInRange(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
		f.parser.EXPECT().Parse(gomock.Any(), gomock.Any()).AnyTimes().Return(juneCalendar(), nil)

		u := uc.New(f.r, f.parser, f.wdCalc, f.salCalc, f.vacPay)
		_, err := u.Do(context.Background(), in)
		assert.Error(t, err)
	})
}
