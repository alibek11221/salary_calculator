package get_salary_report

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	wc "salary_calculator/internal/pkg/http/work_calendar"
	"salary_calculator/internal/pkg/utils"

	eg "golang.org/x/sync/errgroup"
)

type usecase struct {
	r                  repo
	workdaysClient     workdaysClient
	workdaysCalculator workdaysCalculator
	salaryCalculator   salaryCalculator
}

func New(
	r repo,
	workdaysClient workdaysClient,
	workdaysCalculator workdaysCalculator,
	salaryCalculator salaryCalculator,
) *usecase {
	return &usecase{
		r:                  r,
		salaryCalculator:   salaryCalculator,
		workdaysClient:     workdaysClient,
		workdaysCalculator: workdaysCalculator,
	}
}

func (u *usecase) Do(ctx context.Context, in get_salary_report.In) (*get_salary_report.Out, error) {
	targetDate := value_objects.From(in.Year, in.Month)

	var (
		latestSalary *dbstore.SalaryChange
		b            dbstore.Bonuse
		wdr          *wc.WorkdayResponse
	)

	g, gCtx := eg.WithContext(ctx)

	g.Go(func() error {
		var err error
		latestSalary, err = u.getLatestChange(gCtx, targetDate)

		return err
	})

	g.Go(func() error {
		var err error
		b, err = u.r.GetBonusByDate(gCtx, targetDate.String())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		return nil
	})

	g.Go(func() error {
		var err error
		wdr, err = u.workdaysClient.GetWorkdaysForMonth(gCtx, in.Month, in.Year)

		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if latestSalary == nil {
		return nil, fmt.Errorf("latest salary not found for date %s", targetDate.String())
	}

	wDays := u.workdaysCalculator.CalculateWorkDaysForMonth(wdr)
	ndfl := utils.CalculateNDFL(latestSalary.Salary)
	sCtx := value_objects.NewSalaryContext(latestSalary.Salary, ndfl, wDays)

	foodPay := 529 * wDays.TotalWorkdays
	extraCollection := *value_objects.NewExtraPaymentsCollection(value_objects.ExtraPayment{
		Value: float64(foodPay),
		Name:  "За еду",
		T:     value_objects.Salary,
	})

	if b.ID.Valid {
		extraCollection.Push(value_objects.ExtraPayment{
			Value: b.Value,
			Name:  "Бонус",
			T:     value_objects.Advance,
		})
	}

	calc := u.salaryCalculator.CalculateSalary(sCtx, extraCollection)

	return &get_salary_report.Out{
		BaseSalary: latestSalary.Salary,
		Result:     calc,
	}, nil
}

func (u *usecase) getLatestChange(
	ctx context.Context,
	targetDate *value_objects.SalaryDate,
) (
	*dbstore.SalaryChange,
	error,
) {
	changes, err := u.r.ListChanges(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list salary changes: %w", err)
	}

	var latestChangeBeforeTarget *dbstore.SalaryChange
	var latestDateBeforeTarget *value_objects.SalaryDate

	for _, change := range changes {
		changeDate, err := value_objects.NewSalaryDate(change.ChangeFrom)
		if err != nil {
			continue
		}

		if changeDate.Compare(targetDate) <= 0 {
			if latestDateBeforeTarget == nil || changeDate.Compare(latestDateBeforeTarget) > 0 {
				latestChangeBeforeTarget = &change
				latestDateBeforeTarget = changeDate
			}
		}
	}

	if latestChangeBeforeTarget == nil {
		return nil, errors.New("failed to find any salary change before target date")
	}

	return latestChangeBeforeTarget, nil
}
