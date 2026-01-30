package work_calendar

import (
	"salary_calculator/internal/pkg/types"
)

type WorkdayResponse struct {
	Statistics Statistics `json:"statistic"`
	Days       []Day      `json:"days"`
}

type Statistics struct {
	WorkDays             int `json:"work_days"`
	CalendarDays         int `json:"calendar_days"`
	WorkingHours         int `json:"working_hours"`
	ShortenedWorkingDays int `json:"shortened_working_days"`
}

type Day struct {
	Date     types.Date `json:"date"`
	TypeId   int        `json:"type_id"`
	TypeText string     `json:"type_text"`
	WeekDay  string     `json:"week_day"`
}
