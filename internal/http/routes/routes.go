package routes

import (
	"salary_calculator/internal/app"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
)

type Registrar struct {
	registrars []RouterRegistrar
}

func NewRoutesRegistrar(a *app.App) *Registrar {
	registrars := []RouterRegistrar{
		NewHealthRoutesRegistrar(a),
		NewSalaryRoutesRegistrar(a),
		NewBonusRoutesRegistrar(a),
		NewDutyRoutesRegistrar(a),
	}

	return &Registrar{
		registrars: registrars,
	}
}

func (rr *Registrar) RegisterAll(router chi.Router) {
	router.Route("/api/v1", func(r chi.Router) {
		for _, registrar := range rr.registrars {
			r.Route(spew.Sprintf("/%s", registrar.Name()), func(r chi.Router) {
				registrar.Register(r)
			})
		}
	})
}

func (rr *Registrar) GetRegistrars() []RouterRegistrar {
	return rr.registrars
}
