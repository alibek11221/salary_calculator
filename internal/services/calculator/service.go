package calculator

import (
	"context"
	"database/sql"
	"errors"

	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/utils"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	r repo
}

func New(r repo) *Service {
	return &Service{r: r}
}

const (
	FoodPaymentForDay = 529
	FoodPaymentName   = "За еду"

	WorkDayDutyHours = 12
	HolidayDutyHours = 24
	DutyPaymentName  = "За дежурство"

	BonusPaymentName = "Премия"
)

// PayoutBreakdown — разбивка сумм по выплатам: аванс, остаток и их сумма.
type PayoutBreakdown struct {
	Advance float64 `json:"advance"`
	Salary  float64 `json:"salary"`
	Total   float64 `json:"total"`
}

type SalaryCalculationResult struct {
	// Brutto — оклад грязными, без доплат.
	Brutto PayoutBreakdown `json:"brutto"`
	// Netto — оклад за вычетом НДФЛ, без доплат.
	Netto PayoutBreakdown `json:"netto"`
	// InHand — фактические выплаты на руки: netto + доплаты своего типа,
	// Total дополнительно включает выплаты типа extra (премия, отпускные).
	InHand        PayoutBreakdown                          `json:"in_hand"`
	ExtraPayments value_objects.ExtraPaymentsCollectionDto `json:"extra_payments"`
}

func (s *Service) calculateExtraPayments(
	ctx context.Context,
	date value_objects.SalaryDate,
	sCtx value_objects.SalaryCalculationContext,
) (*value_objects.ExtraPaymentsCollection, error) {
	var (
		bonus dbstore.Bonuse
		duty  dbstore.Duty
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		b, err := s.r.GetBonusByDate(gCtx, date.String())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		bonus = b
		return nil
	})

	g.Go(func() error {
		d, err := s.r.GetDutyByDate(gCtx, date.PreviousMonth().String())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		duty = d
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return s.buildExtraPaymentsCollection(date, sCtx, bonus, duty), nil
}

func (s *Service) buildExtraPaymentsCollection(
	date value_objects.SalaryDate,
	sCtx value_objects.SalaryCalculationContext,
	bonus dbstore.Bonuse,
	duty dbstore.Duty,
) *value_objects.ExtraPaymentsCollection {
	foodPayment := s.calculateFoodPayment(sCtx)
	extraCollection := value_objects.NewExtraPaymentsCollection(foodPayment)

	if bonus.ID.Valid {
		bonusPayment := s.calculateBonusPayment(bonus, sCtx)
		extraCollection.Push(bonusPayment)
	}

	if duty.ID.Valid {
		dutyPayment := s.calculateDutyPayment(date, sCtx, duty)
		extraCollection.Push(dutyPayment)
	}

	for _, vp := range sCtx.VacationPayments() {
		extraCollection.Push(vp)
	}

	return extraCollection
}

func (s *Service) calculateFoodPayment(sCtx value_objects.SalaryCalculationContext) value_objects.ExtraPayment {
	// Еда платится за фактически отработанные дни: половинки месяца уже
	// уменьшены на дни отпуска, в обычном месяце их сумма равна TotalWorkdays.
	attendedDays := sCtx.Workdays().FirstHalfDays + sCtx.Workdays().SecondHalfDays
	foodPay := float64(FoodPaymentForDay * attendedDays)
	return value_objects.ExtraPayment{
		Value: utils.ToTwoDecimals(utils.SubPercentage(foodPay, sCtx.CurrentNDFL())),
		Name:  FoodPaymentName,
		T:     value_objects.Salary,
	}
}

func (s *Service) calculateBonusPayment(
	bonus dbstore.Bonuse,
	sCtx value_objects.SalaryCalculationContext,
) value_objects.ExtraPayment {
	return value_objects.ExtraPayment{
		Value: utils.ToTwoDecimals(utils.SubPercentage(bonus.Value, sCtx.CurrentNDFL())),
		Name:  BonusPaymentName,
		T:     value_objects.Extra,
	}
}

func (s *Service) calculateDutyPayment(
	date value_objects.SalaryDate,
	sCtx value_objects.SalaryCalculationContext,
	duty dbstore.Duty,
) value_objects.ExtraPayment {
	lastMonthDays := date.PreviousMonth().CalendarDays()
	hourlyRate := sCtx.CurrentBase() / float64(lastMonthDays) / 24
	totalDutyHours := WorkDayDutyHours*duty.InWorkdays + HolidayDutyHours*duty.InHolidays
	totalDutyPay := hourlyRate * float64(totalDutyHours)
	advancePayment := totalDutyPay / 5

	return value_objects.ExtraPayment{
		Value: utils.ToTwoDecimals(utils.SubPercentage(advancePayment, sCtx.CurrentNDFL())),
		Name:  DutyPaymentName,
		T:     value_objects.Advance,
	}
}

func (s *Service) CalculateSalary(
	ctx context.Context,
	date value_objects.SalaryDate,
	sCtx value_objects.SalaryCalculationContext,
) (*SalaryCalculationResult, error) {
	grossAdvance := utils.ToTwoDecimals(s.calculateGrossAmount(
		sCtx.CurrentBase(), sCtx.Workdays().TotalWorkdays, sCtx.Workdays().FirstHalfDays))
	grossSalary := utils.ToTwoDecimals(s.calculateGrossAmount(
		sCtx.CurrentBase(), sCtx.Workdays().TotalWorkdays, sCtx.Workdays().SecondHalfDays))

	netAdvance := utils.ToTwoDecimals(utils.SubPercentage(grossAdvance, sCtx.CurrentNDFL()))
	netSalary := utils.ToTwoDecimals(utils.SubPercentage(grossSalary, sCtx.CurrentNDFL()))

	extraPayments, err := s.calculateExtraPayments(ctx, date, sCtx)
	if err != nil {
		return nil, err
	}

	totals := extraPayments.Total()
	inHandAdvance := utils.ToTwoDecimals(netAdvance + totals[value_objects.Advance])
	inHandSalary := utils.ToTwoDecimals(netSalary + totals[value_objects.Salary])

	return &SalaryCalculationResult{
		Brutto: PayoutBreakdown{
			Advance: grossAdvance,
			Salary:  grossSalary,
			Total:   utils.ToTwoDecimals(grossAdvance + grossSalary),
		},
		Netto: PayoutBreakdown{
			Advance: netAdvance,
			Salary:  netSalary,
			Total:   utils.ToTwoDecimals(netAdvance + netSalary),
		},
		InHand: PayoutBreakdown{
			Advance: inHandAdvance,
			Salary:  inHandSalary,
			Total:   utils.ToTwoDecimals(inHandAdvance + inHandSalary + totals[value_objects.Extra]),
		},
		ExtraPayments: extraPayments.ToDto(),
	}, nil
}

func (s *Service) calculateGrossAmount(base float64, totalDays, workedDays int) float64 {
	if totalDays <= 0 {
		return 0
	}
	return base / float64(totalDays) * float64(workedDays)
}
