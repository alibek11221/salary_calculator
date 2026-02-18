package value_objects

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const StartYear = 2024

type SalaryDate struct {
	month int
	year  int
	key   string
}

func (s *SalaryDate) String() string {
	if s == nil {
		return ""
	}
	return s.key
}

func parseSalaryDate(date string) (year, month int, err error) {
	splitDate := strings.Split(date, "_")
	if len(splitDate) != 2 {
		return 0, 0, fmt.Errorf("invalid date format: expected 'year_month', got '%s'", date)
	}

	y, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year: %w", err)
	}
	if y < StartYear {
		return 0, 0, fmt.Errorf("year out of range: must be >= %d, got %d", StartYear, y)
	}

	m, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid month: %w", err)
	}

	if m > 12 || m < 1 {
		return 0, 0, fmt.Errorf("month out of range: must be 1-12, got %d", m)
	}

	return y, m, nil
}

func NewSalaryDate(date string) (*SalaryDate, error) {
	y, m, err := parseSalaryDate(date)
	if err != nil {
		return nil, err
	}

	return &SalaryDate{
		year:  y,
		month: m,
		key:   date,
	}, nil
}

func From(year, month int) *SalaryDate {
	return &SalaryDate{year: year, month: month, key: fmt.Sprintf("%d_%02d", year, month)}
}

func (s *SalaryDate) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}

	return json.Marshal(s.key)
}

func (s *SalaryDate) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	y, m, err := parseSalaryDate(str)
	if err != nil {
		return err
	}

	s.year = y
	s.month = m
	s.key = str
	return nil
}

func (s *SalaryDate) Compare(other *SalaryDate) int {
	if s == nil && other == nil {
		return 0
	}
	if s == nil {
		return -1
	}
	if other == nil {
		return 1
	}

	if s.year < other.year {
		return -1
	}
	if s.year > other.year {
		return 1
	}
	if s.month < other.month {
		return -1
	}
	if s.month > other.month {
		return 1
	}
	return 0
}

func (s *SalaryDate) PreviousMonth() *SalaryDate {
	if s.month == 1 {
		return From(s.year-1, 12)
	}

	return From(s.year, s.month-1)
}

func (s *SalaryDate) CalendarDays() int {
	if s == nil {
		return 0
	}

	t := time.Date(s.year, time.Month(s.month+1), 0, 0, 0, 0, 0, time.UTC)
	return t.Day()
}
