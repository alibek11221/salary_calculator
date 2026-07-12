package estimate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	estimate_dto "salary_calculator/internal/dto/estimate_vacation"
	report_dto "salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/utils"
	"salary_calculator/internal/services/vacation_pay"

	"github.com/jackc/pgx/v5/pgtype"
)

type usecase struct {
	r           repo
	vacationPay vacationPayService
	report      reportUsecase
}

func New(r repo, vacationPay vacationPayService, report reportUsecase) *usecase {
	return &usecase{r: r, vacationPay: vacationPay, report: report}
}

func (u *usecase) Do(ctx context.Context, in estimate_dto.In) (*estimate_dto.Out, error) {
	from, to := in.From.Time, in.To.Time

	if err := vacation_pay.ValidatePeriod(from, to); err != nil {
		return nil, err
	}

	// Данные для среднего заработка грузятся один раз на весь estimate-запрос:
	// и для собственного расчёта, и для каждого месячного отчёта ниже.
	earnings, err := u.vacationPay.LoadEarningsData(ctx)
	if err != nil {
		return nil, err
	}

	pay, err := u.vacationPay.CalculatePayWith(earnings, from, to)
	if err != nil {
		return nil, err
	}

	startMonth := value_objects.From(from.Year(), int(from.Month()))
	latest, err := u.r.GetLatestChangeBeforeDate(ctx, startMonth.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no salary data for month %s", startMonth.String())
		}
		return nil, err
	}
	ndfl := utils.CalculateNDFL(latest.Salary)

	hypothetical := dbstore.Vacation{
		DateFrom: pgtype.Date{Time: from, Valid: true},
		DateTo:   pgtype.Date{Time: to, Valid: true},
	}

	var months []estimate_dto.MonthSummary
	for cur := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC); !cur.After(to); cur = cur.AddDate(0, 1, 0) {
		rep, err := u.report.Do(ctx, report_dto.In{
			Year:           cur.Year(),
			Month:          int(cur.Month()),
			ExtraVacations: []dbstore.Vacation{hypothetical},
			Earnings:       earnings,
		})
		if err != nil {
			return nil, err
		}

		months = append(months, estimate_dto.MonthSummary{
			Month:   value_objects.From(cur.Year(), int(cur.Month())).String(),
			Advance: rep.Result.Advance,
			Salary:  rep.Result.Salary,
			Total:   rep.Result.Total,
		})
	}

	return &estimate_dto.Out{
		AvgDaily:       utils.ToTwoDecimals(pay.AvgDaily),
		PaidDays:       pay.PaidDays,
		Gross:          utils.ToTwoDecimals(pay.Gross),
		Net:            utils.ToTwoDecimals(utils.SubPercentage(pay.Gross, ndfl)),
		AffectedMonths: months,
	}, nil
}
