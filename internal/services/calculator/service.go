package calculator

import (
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/pkg/utils"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

type SalaryCalculationResult struct {
	Advance       float64                                  `json:"advance"`
	Salary        float64                                  `json:"salary"`
	Total         float64                                  `json:"total"`
	GrossAdvance  float64                                  `json:"grossAdvance"`
	GrossSalary   float64                                  `json:"grossSalary"`
	GrossTotal    float64                                  `json:"grossTotal"`
	ExtraPayments value_objects.ExtraPaymentsCollectionDto `json:"extra_payments"`
}

func (s *Service) CalculateSalary(
	sCtx value_objects.SalaryCalculationContext,
	extraPayments value_objects.ExtraPaymentsCollection,
) SalaryCalculationResult {
	out := s.calculateResult(sCtx, extraPayments)

	return out
}

func (s *Service) calculateResult(
	sCtx value_objects.SalaryCalculationContext,
	extraPayments value_objects.ExtraPaymentsCollection,
) SalaryCalculationResult {
	res := SalaryCalculationResult{}
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

	extra := extraPayments
	if extra.Type() == value_objects.Advance {
		res.Advance += extra.Total()
	}
	if extra.Type() == value_objects.Salary {
		res.Salary += extra.Total()
	}

	res.ExtraPayments = extra.ToDto()

	res.GrossTotal = res.GrossSalary + res.GrossAdvance
	res.Total = res.Salary + res.Advance

	return res
}

func (s *Service) calculateGrossAmount(base float64, totalDays, workedDays int) float64 {
	out := base / float64(totalDays) * float64(workedDays)

	return out
}
