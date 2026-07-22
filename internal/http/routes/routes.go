package routes

import (
	"salary_calculator/internal/app"
	"salary_calculator/internal/pkg/http/work_calendar_parser"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/vacation_pay"
	"salary_calculator/internal/services/work_days"

	"github.com/go-chi/chi/v5"
)

// SharedServices создаются один раз, чтобы парсер календаря (и его LRU-кэш)
// не дублировался между регистраторами роутов.
type SharedServices struct {
	WorkdaysParser     *work_calendar_parser.Parser
	WorkdaysCalculator *work_days.Service
	SalaryCalculator   *calculator.Service
	VacationPay        *vacation_pay.Service
}

type Registrar struct {
	registrars []RouterRegistrar
}

func NewRoutesRegistrar(a *app.App) *Registrar {
	workdaysParser := work_calendar_parser.New(a.Config.WorkdaysConfig.Dir, a.Config.WorkdaysConfig.CacheCap, a.Logger)
	shared := &SharedServices{
		WorkdaysParser:     workdaysParser,
		WorkdaysCalculator: work_days.New(),
		SalaryCalculator:   calculator.New(a.Repo),
		VacationPay:        vacation_pay.New(a.Repo, workdaysParser),
	}

	registrars := []RouterRegistrar{
		NewHealthRoutesRegistrar(a),
		NewSalaryRoutesRegistrar(a, shared),
		NewBonusRoutesRegistrar(a),
		NewDutyRoutesRegistrar(a),
		NewVacationRoutesRegistrar(a, shared),
	}

	return &Registrar{
		registrars: registrars,
	}
}

func (rr *Registrar) RegisterAll(router chi.Router) {
	router.Route("/api/v1", func(r chi.Router) {
		for _, registrar := range rr.registrars {
			r.Route("/"+registrar.Name(), func(r chi.Router) {
				registrar.Register(r)
			})
		}
	})
}

func (rr *Registrar) GetRegistrars() []RouterRegistrar {
	return rr.registrars
}
