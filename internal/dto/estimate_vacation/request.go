package estimate_vacation

import "salary_calculator/internal/pkg/types"

type In struct {
	From types.ISODate `json:"from"`
	To   types.ISODate `json:"to"`
}

type MonthSummary struct {
	Month   string  `json:"month"`
	Advance float64 `json:"advance"`
	Salary  float64 `json:"salary"`
	Total   float64 `json:"total"`
}

type Out struct {
	AvgDaily       float64        `json:"avg_daily"`
	PaidDays       int            `json:"paid_days"`
	Gross          float64        `json:"gross"`
	Net            float64        `json:"net"`
	AffectedMonths []MonthSummary `json:"affected_months"`
}
