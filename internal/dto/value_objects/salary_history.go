package value_objects

type SalaryHistory interface {
	Add(date SalaryDate, amount float64)
	Get(date SalaryDate) (float64, bool)
}

type salaryHistory struct {
	history map[SalaryDate]float64
}

func NewSalaryHistory() SalaryHistory {
	return &salaryHistory{
		history: make(map[SalaryDate]float64),
	}
}

func (s *salaryHistory) Add(date SalaryDate, amount float64) {
	s.history[date] = amount
}

func (s *salaryHistory) Get(date SalaryDate) (float64, bool) {
	amount, ok := s.history[date]

	return amount, ok
}

func (s *salaryHistory) GetLatestSalary(date *SalaryDate) (*SalaryDate, *float64) {
	var latestDateBeforeTarget *SalaryDate
	var latestValueBeforeTarget *float64
	for changeDate, value := range s.history {
		if changeDate.Compare(date) <= 0 {
			if latestDateBeforeTarget == nil || changeDate.Compare(latestDateBeforeTarget) > 0 {
				latestDateBeforeTarget = &changeDate
				latestValueBeforeTarget = &value
			}
		}
	}

	return latestDateBeforeTarget, latestValueBeforeTarget
}
