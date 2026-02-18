package routes

import (
	"salary_calculator/internal/app"
	health2 "salary_calculator/internal/http/handlers/health"
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
	var checkers []health2.Checker

	checkers = append(checkers, health2.NewBasicHealthChecker())

	if h.app.DB != nil {
		checkers = append(checkers, database.NewHealthChecker(h.app.DB))
	}

	healthHandler := health2.New(checkers, 0)

	router.Get("/health", healthHandler.ServeHTTP)
}

func (h *HealthRoutesRegistrar) Name() string {
	return "health"
}
