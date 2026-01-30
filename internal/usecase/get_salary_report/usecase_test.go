package get_salary_report_test

import (
	"context"
	"database/sql"
	"errors"
	"salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/http/work_calendar"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/work_days"
	getsalaryreportuc "salary_calculator/internal/usecase/get_salary_report"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_Do(t *testing.T) {
	type fields struct {
		r                  *Mockrepo
		workdaysClient     *MockworkdaysClient
		workdaysCalculator *MockworkdaysCalculator
		salaryCalculator   *MocksalaryCalculator
	}

	validIn := get_salary_report.In{
		Year:  2025,
		Month: 1,
	}

	tests := []struct {
		name    string
		setup   func(f fields)
		in      get_salary_report.In
		want    *get_salary_report.Out
		wantErr bool
	}{
		{
			name: "success",
			in:   validIn,
			setup: func(f fields) {
				f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), "2025_01").Return(dbstore.SalaryChange{
					Salary:     100000,
					ChangeFrom: "2025_01",
				}, nil)
				f.r.EXPECT().GetBonusByDate(gomock.Any(), "2025_01").Return(dbstore.Bonuse{}, nil)

				f.workdaysClient.EXPECT().GetWorkdaysForMonth(gomock.Any(), 1, 2025).Return(&work_calendar.WorkdayResponse{}, nil)

				f.workdaysCalculator.EXPECT().CalculateWorkDaysForMonth(gomock.Any()).Return(&work_days.WorkdaysForMonth{
					TotalWorkdays: 20,
				})

				f.salaryCalculator.EXPECT().CalculateSalary(gomock.Any(), gomock.Any()).Return(calculator.SalaryCalculationResult{
					Advance: 50000,
					Salary:  50000,
					Total:   100000,
				})
			},
			want: &get_salary_report.Out{
				BaseSalary: 100000,
				Result: calculator.SalaryCalculationResult{
					Advance: 50000,
					Salary:  50000,
					Total:   100000,
				},
			},
		},
		{
			name: "repo error",
			in:   validIn,
			setup: func(f fields) {
				f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), gomock.Any()).Return(dbstore.SalaryChange{}, errors.New("db error"))
				f.workdaysClient.EXPECT().GetWorkdaysForMonth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				f.r.EXPECT().GetBonusByDate(gomock.Any(), gomock.Any()).AnyTimes().Return(dbstore.Bonuse{}, nil)
			},
			wantErr: true,
		},
		{
			name: "workdays client error",
			in:   validIn,
			setup: func(f fields) {
				f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), gomock.Any()).AnyTimes().Return(dbstore.SalaryChange{Salary: 100000, ChangeFrom: "2025_01"}, nil)
				f.r.EXPECT().GetBonusByDate(gomock.Any(), gomock.Any()).AnyTimes().Return(dbstore.Bonuse{}, nil)
				f.workdaysClient.EXPECT().GetWorkdaysForMonth(gomock.Any(), 1, 2025).Return(nil, errors.New("http error"))
			},
			wantErr: true,
		},
		{
			name: "no salary change found",
			in:   validIn,
			setup: func(f fields) {
				f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), gomock.Any()).Return(dbstore.SalaryChange{}, sql.ErrNoRows)
				f.workdaysClient.EXPECT().GetWorkdaysForMonth(gomock.Any(), 1, 2025).AnyTimes().Return(&work_calendar.WorkdayResponse{}, nil)
				f.r.EXPECT().GetBonusByDate(gomock.Any(), gomock.Any()).AnyTimes().Return(dbstore.Bonuse{}, nil)
			},
			wantErr: true,
		},
		{
			name: "invalid input date",
			in: get_salary_report.In{
				Year:  2023,
				Month: 1,
			},
			setup: func(f fields) {
				f.r.EXPECT().GetLatestChangeBeforeDate(gomock.Any(), "2023_01").Return(dbstore.SalaryChange{}, sql.ErrNoRows)
				f.workdaysClient.EXPECT().GetWorkdaysForMonth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&work_calendar.WorkdayResponse{}, nil)
				f.r.EXPECT().GetBonusByDate(gomock.Any(), gomock.Any()).AnyTimes().Return(dbstore.Bonuse{}, nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				r:                  NewMockrepo(ctrl),
				workdaysClient:     NewMockworkdaysClient(ctrl),
				workdaysCalculator: NewMockworkdaysCalculator(ctrl),
				salaryCalculator:   NewMocksalaryCalculator(ctrl),
			}
			if tt.setup != nil {
				tt.setup(f)
			}

			u := getsalaryreportuc.New(f.r, f.workdaysClient, f.workdaysCalculator, f.salaryCalculator)
			got, err := u.Do(context.Background(), tt.in)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.BaseSalary, got.BaseSalary)
				assert.Equal(t, tt.want.Result.Advance, got.Result.Advance)
			}
		})
	}
}
