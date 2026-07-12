package vacation_pay

import (
	"errors"
	"fmt"
	"time"

	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/http/work_calendar_parser"
)

// AvgCalendarDaysPerMonth — среднемесячное число календарных дней (ст. 139 ТК РФ).
const AvgCalendarDaysPerMonth = 29.3

const (
	dayTypeWorkday   = 1
	dayTypeHoliday   = 3
	dayTypeShortened = 5
)

var (
	ErrInvalidPeriod  = errors.New("invalid vacation period")
	ErrNoEarningsData = errors.New("no salary data for average earnings period")
)

type Service struct {
	r      repo
	parser calendarParser
}

func New(r repo, parser calendarParser) *Service {
	return &Service{r: r, parser: parser}
}

// Deduction — сколько рабочих дней месяца пришлось на отпуск, по половинкам (1–15 / 16+).
type Deduction struct {
	FirstHalf  int
	SecondHalf int
}

// ValidatePeriod проверяет базовую корректность периода отпуска.
func ValidatePeriod(from, to time.Time) error {
	if from.IsZero() || to.IsZero() {
		return ErrInvalidPeriod
	}
	if to.Before(from) {
		return ErrInvalidPeriod
	}
	if from.Year() < value_objects.StartYear {
		return ErrInvalidPeriod
	}
	return nil
}

// PaidDays — оплачиваемые дни отпуска: календарные дни периода минус
// госпраздники (ст. 120 ТК: праздники в отпуск не входят). Выходные входят.
func (s *Service) PaidDays(from, to time.Time) (int, error) {
	if to.Before(from) {
		return 0, ErrInvalidPeriod
	}

	paid := 0
	end := truncateToDay(to)
	for d := truncateToDay(from); !d.After(end); d = d.AddDate(0, 0, 1) {
		cal, err := s.parser.Parse(d.Year(), int(d.Month()))
		if err != nil {
			return 0, err
		}
		typeID, ok := dayTypeFor(cal, d)
		if !ok {
			return 0, fmt.Errorf("day %s not found in work calendar", d.Format("2006-01-02"))
		}
		if typeID != dayTypeHoliday {
			paid++
		}
	}
	return paid, nil
}

// WorkdayDeduction считает, сколько рабочих дней (тип 1 и 5) из переданного
// месячного календаря покрыто отпусками, с разбивкой по половинам месяца.
func WorkdayDeduction(days []work_calendar_parser.Day, vacations []dbstore.Vacation) Deduction {
	var ded Deduction
	for _, d := range days {
		if d.TypeId != dayTypeWorkday && d.TypeId != dayTypeShortened {
			continue
		}
		if !coveredByAny(d.Date.Time, vacations) {
			continue
		}
		if d.Date.Day() <= 15 {
			ded.FirstHalf++
		} else {
			ded.SecondHalf++
		}
	}
	return ded
}

func coveredByAny(d time.Time, vacations []dbstore.Vacation) bool {
	day := truncateToDay(d)
	for _, v := range vacations {
		if !day.Before(truncateToDay(v.DateFrom.Time)) && !day.After(truncateToDay(v.DateTo.Time)) {
			return true
		}
	}
	return false
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func dayTypeFor(cal *work_calendar_parser.WorkdayResponse, d time.Time) (int, bool) {
	if cal == nil {
		return 0, false
	}
	for _, day := range cal.Days {
		if day.Date.Year() == d.Year() && day.Date.Month() == d.Month() && day.Date.Day() == d.Day() {
			return day.TypeId, true
		}
	}
	return 0, false
}
