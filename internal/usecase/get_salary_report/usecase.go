package get_salary_report

import (
	"context"
	"errors"
	"fmt"

	"salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/http/work_calendar"
	"salary_calculator/internal/pkg/utils"

	"golang.org/x/sync/errgroup"
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
	eg, ctx := errgroup.WithContext(ctx)

	var latestSalary *dbstore.SalaryChange
	var wdr *work_calendar.WorkdayResponse

	eg.Go(func() error {
		var err error
		latestSalary, err = u.getLatestChange(ctx, in)
		return err
	})

	eg.Go(func() error {
		var err error
		wdr, err = u.workdaysClient.GetWorkdaysForMonth(ctx, in.Month, in.Year)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	wDays := u.workdaysCalculator.CalculateWorkDaysForMonth(*wdr)
	ndfl := utils.CalculateNDFL(latestSalary.Salary)
	sCtx := value_objects.NewSalaryContext(latestSalary.Salary, ndfl, wDays)

	foodPay := 529 * wDays.TotalWorkdays
	extraCollection := *value_objects.NewExtraPaymentsCollection(value_objects.Salary, value_objects.ExtraPayment{
		Value: float64(foodPay),
		Name:  "За еду",
	})

	calc := u.salaryCalculator.CalculateSalary(sCtx, extraCollection)

	return &get_salary_report.Out{
		BaseSalary: latestSalary.Salary,
		Result:     calc,
	}, nil
}

func (u *usecase) getLatestChange(ctx context.Context, in get_salary_report.In) (*dbstore.SalaryChange, error) {
	changes, err := u.r.ListChanges(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list salary changes: %w", err)
	}

	targetDate, err := value_objects.NewSalaryDate(fmt.Sprintf("%d_%d", in.Year, in.Month))
	if err != nil {
		return nil, fmt.Errorf("invalid input date: %w", err)
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
