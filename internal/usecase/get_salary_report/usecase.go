package get_salary_report

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	wc "salary_calculator/internal/pkg/http/work_calendar_parser"
	"salary_calculator/internal/pkg/utils"
	"salary_calculator/internal/services/vacation_pay"

	"github.com/jackc/pgx/v5/pgtype"
	eg "golang.org/x/sync/errgroup"
)

type usecase struct {
	r                  repo
	workdaysParser     workdaysParser
	workdaysCalculator workdaysCalculator
	salaryCalculator   salaryCalculator
	vacationPay        vacationPay
}

func New(
	r repo,
	workdaysParser workdaysParser,
	workdaysCalculator workdaysCalculator,
	salaryCalculator salaryCalculator,
	vacationPay vacationPay,
) *usecase {
	return &usecase{
		r:                  r,
		salaryCalculator:   salaryCalculator,
		workdaysParser:     workdaysParser,
		workdaysCalculator: workdaysCalculator,
		vacationPay:        vacationPay,
	}
}

func (u *usecase) Do(ctx context.Context, in get_salary_report.In) (*get_salary_report.Out, error) {
	targetDate := value_objects.From(in.Year, in.Month)

	var (
		latestSalary *dbstore.SalaryChange
		wdr          *wc.WorkdayResponse
		dbVacations  []dbstore.Vacation
	)

	monthStart := time.Date(in.Year, time.Month(in.Month), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, -1)

	g, gCtx := eg.WithContext(ctx)

	g.Go(func() error {
		latestSalaryRow, err := u.r.GetLatestChangeBeforeDate(gCtx, targetDate.String())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("latest salary not found for date %s", targetDate.String())
			}
			return err
		}
		latestSalary = &latestSalaryRow

		return nil
	})

	g.Go(func() error {
		var err error
		wdr, err = u.workdaysParser.Parse(in.Year, in.Month)

		return err
	})

	g.Go(func() error {
		var err error
		dbVacations, err = u.r.GetVacationsInRange(gCtx, dbstore.GetVacationsInRangeParams{
			RangeStart: pgtype.Date{Time: monthStart, Valid: true},
			RangeEnd:   pgtype.Date{Time: monthEnd, Valid: true},
		})

		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if latestSalary == nil {
		return nil, fmt.Errorf("latest salary not found for date %s", targetDate.String())
	}

	wDays := u.workdaysCalculator.CalculateWorkDaysForMonth(wdr)
	if wDays == nil {
		return nil, fmt.Errorf("could not calculate workdays")
	}

	vacations := make([]dbstore.Vacation, 0, len(dbVacations)+len(in.ExtraVacations))
	vacations = append(vacations, dbVacations...)
	vacations = append(vacations, in.ExtraVacations...)

	vacation_pay.WorkdayDeduction(wdr.Days, vacations).ApplyTo(wDays)

	ndfl := utils.CalculateNDFL(latestSalary.Salary)

	vacationPayments, err := u.vacationPay.NetPaymentsForMonth(ctx, vacations, in.Year, in.Month, ndfl, in.Earnings)
	if err != nil {
		return nil, err
	}

	sCtx := value_objects.NewSalaryContext(latestSalary.Salary, ndfl, *wDays).
		WithVacationPayments(vacationPayments...)

	calc, err := u.salaryCalculator.CalculateSalary(ctx, *targetDate, sCtx)
	if err != nil {
		return nil, err
	}

	return &get_salary_report.Out{
		BaseSalary: latestSalary.Salary,
		Result:     calc,
	}, nil
}
