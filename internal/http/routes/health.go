package routes

import (
	"salary_calculator/internal/app"
	"salary_calculator/internal/http/handlers/service_h/health"
	"salary_calculator/internal/pkg/database"

	"github.com/go-chi/chi/v5"
)

type HealthRoutesRegistrar struct {
	app *app.App
}

func NewHealthRoutesRegistrar(a *app.App) *HealthRoutesRegistrar {
	return &HealthRoutesRegistrar{app: a}
}

func (h *HealthRoutesRegistrar) Register(router chi.Router) {
	var checkers []health.Checker

	checkers = append(checkers, health.NewBasicHealthChecker())

	if h.app.DB != nil {
		checkers = append(checkers, database.NewHealthChecker(h.app.DB))
	}

	healthHandler := health.New(checkers, 0)

	router.Get("/health", healthHandler.ServeHTTP)
}

func (h *HealthRoutesRegistrar) Name() string {
	return "health"
}
