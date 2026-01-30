package server

import (
	"net/http"
	"time"

	_ "salary_calculator/docs"
	"salary_calculator/internal/app"
	"salary_calculator/internal/http/routes"
	_ "salary_calculator/internal/pkg/logging"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewServer(a *app.App) (*http.Server, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByRealIP(100, time.Minute))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	routeRegistrar := routes.NewRoutesRegistrar(a)
	routeRegistrar.RegisterAll(r)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Путь к сгенерированному swagger.json
	))

	server := &http.Server{
		Addr:           ":" + a.Config.Port,
		Handler:        r,
		ReadTimeout:    a.Config.Server.ReadTimeout,
		WriteTimeout:   a.Config.Server.WriteTimeout,
		MaxHeaderBytes: a.Config.Server.MaxHeaderBytes,
	}

	log.Info().
		Str("port", a.Config.Port).
		Msg("server initialized")

	return server, nil
}
