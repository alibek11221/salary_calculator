package calculator_test

import (
	"context"
	"database/sql"
	"testing"

	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/work_days"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func validUUID() pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
}

func TestService_CalculateSalary(t *testing.T) {
	// Февраль 2025: бонус ищется по "2025_02", дежурство — по прошлому месяцу "2025_01" (31 календарный день).
	date := value_objects.From(2025, 2)

	sCtx := value_objects.NewSalaryContext(
		200000,
		13,
		work_days.WorkdaysForMonth{TotalWorkdays: 20, FirstHalfDays: 10, SecondHalfDays: 10},
	)

	tests := []struct {
		name  string
		sCtx  *value_objects.SalaryCalculationContext
		setup func(r *Mockrepo)
		check func(t *testing.T, got *calculator.SalaryCalculationResult)
	}{
		{
			name: "regular month: only food extra",
			setup: func(r *Mockrepo) {
				r.EXPECT().GetBonusByDate(gomock.Any(), "2025_02").Return(dbstore.Bonuse{}, sql.ErrNoRows)
				r.EXPECT().GetDutyByDate(gomock.Any(), "2025_01").Return(dbstore.Duty{}, sql.ErrNoRows)
			},
			check: func(t *testing.T, got *calculator.SalaryCalculationResult) {
				assert.Equal(t, 100000.0, got.GrossAdvance)
				assert.Equal(t, 100000.0, got.GrossSalary)
				assert.Equal(t, 200000.0, got.GrossTotal)
				assert.Equal(t, 87000.0, got.Advance)
				// 87000 + еда 529×20 = 10580
				assert.Equal(t, 97580.0, got.Salary)
				assert.Equal(t, 184580.0, got.Total)
				assert.Len(t, got.ExtraPayments.Payments, 1)
				assert.Equal(t, calculator.FoodPaymentName, got.ExtraPayments.Payments[0].Name)
				assert.Equal(t, 10580.0, got.ExtraPayments.Total)
			},
		},
		{
			name: "month with bonus",
			setup: func(r *Mockrepo) {
				r.EXPECT().GetBonusByDate(gomock.Any(), "2025_02").Return(dbstore.Bonuse{
					ID: validUUID(), Value: 15000, Date: "2025_02",
				}, nil)
				r.EXPECT().GetDutyByDate(gomock.Any(), "2025_01").Return(dbstore.Duty{}, sql.ErrNoRows)
			},
			check: func(t *testing.T, got *calculator.SalaryCalculationResult) {
				// Премия имеет тип extra: в аванс/зарплату не входит, только в списке доплат.
				assert.Equal(t, 87000.0, got.Advance)
				assert.Equal(t, 97580.0, got.Salary)
				assert.Equal(t, 184580.0, got.Total)
				assert.Len(t, got.ExtraPayments.Payments, 2)
				assert.Equal(t, 25580.0, got.ExtraPayments.Total) // 10580 еда + 15000 премия
			},
		},
		{
			name: "month with duty",
			setup: func(r *Mockrepo) {
				r.EXPECT().GetBonusByDate(gomock.Any(), "2025_02").Return(dbstore.Bonuse{}, sql.ErrNoRows)
				r.EXPECT().GetDutyByDate(gomock.Any(), "2025_01").Return(dbstore.Duty{
					ID: validUUID(), Date: "2025_01", InWorkdays: 2, InHolidays: 1,
				}, nil)
			},
			check: func(t *testing.T, got *calculator.SalaryCalculationResult) {
				// 12ч×2 + 24ч×1 = 48ч; ставка 200000/31/24; выплата = ставка×48/5 (аванс).
				expectedDuty := 200000.0 / 31.0 / 24.0 * 48.0 / 5.0
				assert.InDelta(t, 87000.0+expectedDuty, got.Advance, 0.01)
				assert.Equal(t, 97580.0, got.Salary)
				assert.Len(t, got.ExtraPayments.Payments, 2)
			},
		},
		{
			name: "month with vacation",
			sCtx: func() *value_objects.SalaryCalculationContext {
				// Отпуск съел 3 рабочих дня первой половины: 10-3=7. Total остаётся 20.
				c := value_objects.NewSalaryContext(
					200000,
					13,
					work_days.WorkdaysForMonth{TotalWorkdays: 20, FirstHalfDays: 7, SecondHalfDays: 10},
				).WithVacationPayments(value_objects.ExtraPayment{
					Name:  calculator.VacationPaymentName,
					Value: 50000, // net, посчитан usecase'ом
					T:     value_objects.Extra,
				})
				return &c
			}(),
			setup: func(r *Mockrepo) {
				r.EXPECT().GetBonusByDate(gomock.Any(), "2025_02").Return(dbstore.Bonuse{}, sql.ErrNoRows)
				r.EXPECT().GetDutyByDate(gomock.Any(), "2025_01").Return(dbstore.Duty{}, sql.ErrNoRows)
			},
			check: func(t *testing.T, got *calculator.SalaryCalculationResult) {
				// Аванс: 200000/20×7 = 70000 gross → 60900 net.
				assert.Equal(t, 70000.0, got.GrossAdvance)
				assert.Equal(t, 60900.0, got.Advance)
				// Еда по отработанным дням: 529×(7+10)=8993. Зарплата: 87000+8993=95993.
				assert.Equal(t, 95993.0, got.Salary)
				// Строки доплат: еда + отпускные.
				assert.Len(t, got.ExtraPayments.Payments, 2)
				names := []string{got.ExtraPayments.Payments[0].Name, got.ExtraPayments.Payments[1].Name}
				assert.Contains(t, names, calculator.VacationPaymentName)
				assert.Equal(t, 58993.0, got.ExtraPayments.Total) // 8993 + 50000
			},
		},
		{
			name: "repo error propagates",
			setup: func(r *Mockrepo) {
				r.EXPECT().GetBonusByDate(gomock.Any(), gomock.Any()).Return(dbstore.Bonuse{}, assert.AnError)
				r.EXPECT().GetDutyByDate(gomock.Any(), gomock.Any()).AnyTimes().Return(dbstore.Duty{}, sql.ErrNoRows)
			},
			check: nil, // ожидаем ошибку
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockrepo(ctrl)
			tt.setup(r)

			ctxToUse := sCtx
			if tt.sCtx != nil {
				ctxToUse = *tt.sCtx
			}

			svc := calculator.New(r)
			got, err := svc.CalculateSalary(context.Background(), *date, ctxToUse)

			if tt.check == nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			tt.check(t, got)
		})
	}
}
