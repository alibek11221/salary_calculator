package value_objects

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const StartYear = 2024

type SalaryDate struct {
	month int
	year  int
	key   string
}

func (s *SalaryDate) String() string {
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

func (s SalaryDate) MarshalJSON() ([]byte, error) {
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
