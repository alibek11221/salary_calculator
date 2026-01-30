package work_days

import (
	"salary_calculator/internal/pkg/http/work_calendar"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

type WorkdaysForMonth struct {
	TotalWorkdays  int `json:"total_workdays"`
	FirstHalfDays  int `json:"first_half_days"`
	SecondHalfDays int `json:"second_half_days"`
}

func (s *Service) CalculateWorkDays(months map[int]work_calendar.WorkdayResponse) map[int]WorkdaysForMonth {
	out := make(map[int]WorkdaysForMonth, len(months))

	for k, v := range months {
		d := s.CalculateWorkDaysForMonth(v)

		out[k] = d
	}

	return out
}

func (s *Service) CalculateWorkDaysForMonth(month work_calendar.WorkdayResponse) WorkdaysForMonth {
	out := WorkdaysForMonth{
		TotalWorkdays: month.Statistics.WorkDays,
	}
	for _, day := range month.Days {
		if day.TypeId != 1 && day.TypeId != 5 {
			continue
		}

		if day.Date.Day() <= 15 {
			out.FirstHalfDays += 1
			continue
		}
		out.SecondHalfDays += 1
	}

	return out
}
