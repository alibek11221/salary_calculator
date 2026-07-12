package vacation_pay

import (
	"context"

	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/http/work_calendar_parser"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type repo interface {
	ListChanges(ctx context.Context) ([]dbstore.SalaryChange, error)
	ListBonuses(ctx context.Context) ([]dbstore.Bonuse, error)
}

type calendarParser interface {
	Parse(year, month int) (*work_calendar_parser.WorkdayResponse, error)
}
