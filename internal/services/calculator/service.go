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

type SalaryCalculationResult struct {
	Advance       float64                                  `json:"advance"`
	Salary        float64                                  `json:"salary"`
	Total         float64                                  `json:"total"`
	GrossAdvance  float64                                  `json:"grossAdvance"`
	GrossSalary   float64                                  `json:"grossSalary"`
	GrossTotal    float64                                  `json:"grossTotal"`
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
		bonusPayment := s.calculateBonusPayment(bonus)
		extraCollection.Push(bonusPayment)
	}

	if duty.ID.Valid {
		dutyPayment := s.calculateDutyPayment(date, sCtx, duty)
		extraCollection.Push(dutyPayment)
	}

	return extraCollection
}

func (s *Service) calculateFoodPayment(sCtx value_objects.SalaryCalculationContext) value_objects.ExtraPayment {
	foodPay := FoodPaymentForDay * sCtx.Workdays().TotalWorkdays
	return value_objects.ExtraPayment{
		Value: float64(foodPay),
		Name:  FoodPaymentName,
		T:     value_objects.Salary,
	}
}

func (s *Service) calculateBonusPayment(bonus dbstore.Bonuse) value_objects.ExtraPayment {
	return value_objects.ExtraPayment{
		Value: bonus.Value,
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
		Value: advancePayment,
		Name:  DutyPaymentName,
		T:     value_objects.Advance,
	}
}

func (s *Service) CalculateSalary(
	ctx context.Context,
	date value_objects.SalaryDate,
	sCtx value_objects.SalaryCalculationContext,
) (*SalaryCalculationResult, error) {
	res := &SalaryCalculationResult{}
	res.GrossAdvance = utils.ToTwoDecimals(
		s.calculateGrossAmount(
			sCtx.CurrentBase(),
			sCtx.Workdays().TotalWorkdays,
			sCtx.Workdays().FirstHalfDays),
	)
	res.GrossSalary = utils.ToTwoDecimals(
		s.calculateGrossAmount(
			sCtx.CurrentBase(),
			sCtx.Workdays().TotalWorkdays,
			sCtx.Workdays().SecondHalfDays),
	)
	res.Advance = utils.ToTwoDecimals(
		utils.SubPercentage(res.GrossAdvance, sCtx.CurrentNDFL()),
	)
	res.Salary = utils.ToTwoDecimals(
		utils.SubPercentage(res.GrossSalary, sCtx.CurrentNDFL()),
	)

	extraPayments, err := s.calculateExtraPayments(ctx, date, sCtx)
	if err != nil {
		return nil, err
	}

	if extraPayments == nil {
		return res, nil
	}

	res.Advance += extraPayments.Total()[value_objects.Advance]
	res.Salary += extraPayments.Total()[value_objects.Salary]

	res.ExtraPayments = extraPayments.ToDto()

	res.GrossTotal = res.GrossSalary + res.GrossAdvance
	res.Total = res.Salary + res.Advance

	return res, nil
}

func (s *Service) calculateGrossAmount(base float64, totalDays, workedDays int) float64 {
	if totalDays <= 0 {
		return 0
	}
	return base / float64(totalDays) * float64(workedDays)
}
