package calculator

import (
	"testing"

	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/services/work_days"

	"github.com/stretchr/testify/assert"
)

func TestService_CalculateSalary(t *testing.T) {
	service := New()

	tests := []struct {
		name          string
		sCtx          value_objects.SalaryCalculationContext
		extraPayments *value_objects.ExtraPaymentsCollection
		want          SalaryCalculationResult
	}{
		{
			name: "Regular month no extra",
			sCtx: value_objects.NewSalaryContext(
				200000,
				13,
				work_days.WorkdaysForMonth{
					TotalWorkdays:  20,
					FirstHalfDays:  10,
					SecondHalfDays: 10,
				},
			),
			extraPayments: value_objects.NewExtraPaymentsCollection(),
			want: SalaryCalculationResult{
				GrossAdvance:  100000,
				GrossSalary:   100000,
				GrossTotal:    200000,
				Advance:       87000,
				Salary:        87000,
				Total:         174000,
				ExtraPayments: value_objects.ExtraPaymentsCollectionDto{Payments: nil, Total: 0},
			},
		},
		{
			name: "Month with extra in advance",
			sCtx: value_objects.NewSalaryContext(
				200000,
				13,
				work_days.WorkdaysForMonth{
					TotalWorkdays:  20,
					FirstHalfDays:  10,
					SecondHalfDays: 10,
				},
			),
			extraPayments: value_objects.NewExtraPaymentsCollection(value_objects.ExtraPayment{
				Name:  "Bonus",
				Value: 10000,
				T:     value_objects.Advance,
			}),
			want: SalaryCalculationResult{
				GrossAdvance: 100000,
				GrossSalary:  100000,
				GrossTotal:   200000,
				Advance:      97000,
				Salary:       87000,
				Total:        184000,
				ExtraPayments: value_objects.ExtraPaymentsCollectionDto{
					Payments: []value_objects.ExtraPayment{
						{Name: "Bonus", Value: 10000, T: value_objects.Advance},
					},
					Total: 10000,
				},
			},
		},
		{
			name: "Month with extra in salary",
			sCtx: value_objects.NewSalaryContext(
				200000,
				13,
				work_days.WorkdaysForMonth{
					TotalWorkdays:  20,
					FirstHalfDays:  10,
					SecondHalfDays: 10,
				},
			),
			extraPayments: value_objects.NewExtraPaymentsCollection(value_objects.ExtraPayment{
				Name:  "Bonus",
				Value: 10000,
				T:     value_objects.Salary,
			}),
			want: SalaryCalculationResult{
				GrossAdvance: 100000,
				GrossSalary:  100000,
				GrossTotal:   200000,
				Advance:      87000,
				Salary:       97000,
				Total:        184000,
				ExtraPayments: value_objects.ExtraPaymentsCollectionDto{
					Payments: []value_objects.ExtraPayment{
						{Name: "Bonus", Value: 10000, T: value_objects.Salary},
					},
					Total: 10000,
				},
			},
		},
		{
			name: "Uneven days",
			sCtx: value_objects.NewSalaryContext(
				200000,
				13,
				work_days.WorkdaysForMonth{
					TotalWorkdays:  22,
					FirstHalfDays:  10,
					SecondHalfDays: 12,
				},
			),
			extraPayments: value_objects.NewExtraPaymentsCollection(),
			want: SalaryCalculationResult{
				GrossAdvance:  90909.09,
				GrossSalary:   109090.91,
				GrossTotal:    200000,
				Advance:       79090.91,
				Salary:        94909.09,
				Total:         174000,
				ExtraPayments: value_objects.ExtraPaymentsCollectionDto{Payments: nil, Total: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.CalculateSalary(tt.sCtx, *tt.extraPayments)
			assert.Equal(t, tt.want, got)
		})
	}
}
