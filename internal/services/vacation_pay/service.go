package vacation_pay

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/http/work_calendar_parser"
	"salary_calculator/internal/pkg/utils"
	"salary_calculator/internal/services/work_days"
)

// AvgCalendarDaysPerMonth — среднемесячное число календарных дней (ст. 139 ТК РФ).
const AvgCalendarDaysPerMonth = 29.3

const PaymentName = "Отпускные"

// monthsInAvgPeriod — расчётный период среднего заработка (ст. 139 ТК).
const monthsInAvgPeriod = 12

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

func (d Deduction) ApplyTo(w *work_days.WorkdaysForMonth) {
	w.FirstHalfDays -= d.FirstHalf
	w.SecondHalfDays -= d.SecondHalf
}

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

// Pay — результат расчёта отпускных (до НДФЛ).
type Pay struct {
	AvgDaily float64
	PaidDays int
	Gross    float64
}

// EarningsData загружается один раз (LoadEarningsData) и переиспользуется
// между расчётами среднего заработка.
type EarningsData struct {
	Changes []dbstore.SalaryChange
	Bonuses []dbstore.Bonuse
}

// LoadEarningsData — единственное место, где для расчёта среднего заработка
// читаются изменения оклада и премии.
func (s *Service) LoadEarningsData(ctx context.Context) (*EarningsData, error) {
	changes, err := s.r.ListChanges(ctx)
	if err != nil {
		return nil, err
	}
	bonuses, err := s.r.ListBonuses(ctx)
	if err != nil {
		return nil, err
	}
	return &EarningsData{Changes: changes, Bonuses: bonuses}, nil
}

// AverageDailyEarnings — средний дневной заработок за 12 календарных месяцев
// до месяца начала отпуска: (оклад + премии) по месяцам / N / 29.3.
// Месяцы без данных об окладе выпадают и из суммы, и из делителя.
func (s *Service) AverageDailyEarnings(data *EarningsData, vacationStart time.Time) (float64, error) {
	if data == nil {
		return 0, ErrNoEarningsData
	}

	bonusByMonth := make(map[string]float64, len(data.Bonuses))
	for _, b := range data.Bonuses {
		bonusByMonth[b.Date] += b.Value
	}

	timeline := newSalaryTimeline(data.Changes)

	var (
		total          float64
		monthsWithData int
	)

	month := value_objects.From(vacationStart.Year(), int(vacationStart.Month()))
	for range monthsInAvgPeriod {
		month = month.PreviousMonth()
		key := month.String()

		salary, ok := timeline.salaryFor(key)
		if !ok {
			continue
		}

		total += salary + bonusByMonth[key]
		monthsWithData++
	}

	if monthsWithData == 0 {
		return 0, ErrNoEarningsData
	}

	return total / float64(monthsWithData) / AvgCalendarDaysPerMonth, nil
}

// CalculatePayWith — отпускные (gross) за период по уже загруженным данным.
func (s *Service) CalculatePayWith(data *EarningsData, from, to time.Time) (*Pay, error) {
	avg, err := s.AverageDailyEarnings(data, from)
	if err != nil {
		return nil, err
	}

	days, err := s.PaidDays(from, to)
	if err != nil {
		return nil, err
	}

	return &Pay{
		AvgDaily: avg,
		PaidDays: days,
		Gross:    avg * float64(days),
	}, nil
}

// NetPaymentsForMonth пропускает отпуска, начавшиеся вне year/month:
// их выплата целиком показана в месяце начала отпуска.
func (s *Service) NetPaymentsForMonth(
	ctx context.Context,
	vacations []dbstore.Vacation,
	year, month int,
	ndfl float64,
	earnings *EarningsData,
) ([]value_objects.ExtraPayment, error) {
	payments := make([]value_objects.ExtraPayment, 0, len(vacations))

	for _, v := range vacations {
		start := v.DateFrom.Time
		if start.Year() != year || int(start.Month()) != month {
			continue
		}

		if earnings == nil {
			var err error
			earnings, err = s.LoadEarningsData(ctx)
			if err != nil {
				return nil, err
			}
		}

		pay, err := s.CalculatePayWith(earnings, v.DateFrom.Time, v.DateTo.Time)
		if err != nil {
			return nil, err
		}

		payments = append(payments, value_objects.ExtraPayment{
			Name:  PaymentName,
			Value: utils.ToTwoDecimals(utils.SubPercentage(pay.Gross, ndfl)),
			T:     value_objects.Extra,
		})
	}

	return payments, nil
}

func (s *Service) CalculatePay(ctx context.Context, from, to time.Time) (*Pay, error) {
	data, err := s.LoadEarningsData(ctx)
	if err != nil {
		return nil, err
	}
	return s.CalculatePayWith(data, from, to)
}

// salaryTimeline — оклады, отсортированные по месяцу вступления в силу,
// для поиска «последний change_from <= месяца» за O(log n).
// При дубликатах change_from сохраняется первое вхождение (как в исходном
// линейном сканировании).
type salaryTimeline struct {
	keys     []string // отсортированные уникальные change_from ("YYYY_MM")
	salaries map[string]float64
}

func newSalaryTimeline(changes []dbstore.SalaryChange) salaryTimeline {
	salaries := make(map[string]float64, len(changes))
	keys := make([]string, 0, len(changes))
	for _, c := range changes {
		if _, ok := salaries[c.ChangeFrom]; ok {
			continue
		}
		salaries[c.ChangeFrom] = c.Salary
		keys = append(keys, c.ChangeFrom)
	}
	slices.Sort(keys)
	return salaryTimeline{keys: keys, salaries: salaries}
}

// salaryFor — оклад, действующий в месяце "YYYY_MM": последний change_from <= месяца.
// Ключи zero-padded, поэтому лексикографическое сравнение совпадает с хронологией.
func (t salaryTimeline) salaryFor(month string) (float64, bool) {
	i, found := slices.BinarySearch(t.keys, month)
	if found {
		return t.salaries[month], true
	}
	if i == 0 {
		return 0, false
	}
	return t.salaries[t.keys[i-1]], true
}

// PaidDays — оплачиваемые дни отпуска: календарные дни периода минус
// госпраздники (ст. 120 ТК: праздники в отпуск не входят). Выходные входят.
func (s *Service) PaidDays(from, to time.Time) (int, error) {
	if to.Before(from) {
		return 0, ErrInvalidPeriod
	}

	start := truncateToDay(from)
	end := truncateToDay(to)

	paid := 0
	for monthStart := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC); !monthStart.After(end); monthStart = monthStart.AddDate(0, 1, 0) {
		cal, err := s.parser.Parse(monthStart.Year(), int(monthStart.Month()))
		if err != nil {
			return 0, err
		}

		typeByDay := dayTypesFor(cal, monthStart.Year(), monthStart.Month())

		first := monthStart
		if first.Before(start) {
			first = start
		}
		last := monthStart.AddDate(0, 1, -1)
		if last.After(end) {
			last = end
		}

		for d := first.Day(); d <= last.Day(); d++ {
			typeID, ok := typeByDay[d]
			if !ok {
				missing := time.Date(monthStart.Year(), monthStart.Month(), d, 0, 0, 0, 0, time.UTC)
				return 0, fmt.Errorf("day %s not found in work calendar", missing.Format("2006-01-02"))
			}
			if typeID != dayTypeHoliday {
				paid++
			}
		}
	}
	return paid, nil
}

// dayTypesFor — map[день месяца]typeID из месячного календаря;
// дни других месяцев/лет игнорируются.
func dayTypesFor(cal *work_calendar_parser.WorkdayResponse, year int, month time.Month) map[int]int {
	if cal == nil {
		return nil
	}
	types := make(map[int]int, len(cal.Days))
	for _, day := range cal.Days {
		if day.Date.Year() == year && day.Date.Month() == month {
			types[day.Date.Day()] = day.TypeId
		}
	}
	return types
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
