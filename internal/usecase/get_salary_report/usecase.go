package get_salary_report

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	wc "salary_calculator/internal/pkg/http/work_calendar_parser"
	"salary_calculator/internal/pkg/utils"

	eg "golang.org/x/sync/errgroup"
)

type usecase struct {
	r                  repo
	workdaysParser     workdaysParser
	workdaysCalculator workdaysCalculator
	salaryCalculator   salaryCalculator
}

func New(
	r repo,
	workdaysParser workdaysParser,
	workdaysCalculator workdaysCalculator,
	salaryCalculator salaryCalculator,
) *usecase {
	return &usecase{
		r:                  r,
		salaryCalculator:   salaryCalculator,
		workdaysParser:     workdaysParser,
		workdaysCalculator: workdaysCalculator,
	}
}

func (u *usecase) Do(ctx context.Context, in get_salary_report.In) (*get_salary_report.Out, error) {
	targetDate := value_objects.From(in.Year, in.Month)
	var (
		latestSalary *dbstore.SalaryChange
		wdr          *wc.WorkdayResponse
	)

	g, gCtx := eg.WithContext(ctx)

	g.Go(func() error {
		var err error
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

	ndfl := utils.CalculateNDFL(latestSalary.Salary)
	sCtx := value_objects.NewSalaryContext(latestSalary.Salary, ndfl, *wDays)

	calc, err := u.salaryCalculator.CalculateSalary(ctx, *targetDate, sCtx)
	if err != nil {
		return nil, err
	}

	return &get_salary_report.Out{
		BaseSalary: latestSalary.Salary,
		Result:     calc,
	}, nil
}
