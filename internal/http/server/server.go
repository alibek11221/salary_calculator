package server

import (
	"net/http"
	"time"

	_ "salary_calculator/docs"
	"salary_calculator/internal/app"
	"salary_calculator/internal/http/routes"
	"salary_calculator/internal/pkg/logging"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewServer(a *app.App) (*http.Server, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logging.GetChiMiddleware(a.Logger))
	r.Use(middleware.Recoverer)
	// Инвариант: request-timeout < WriteTimeout, чтобы middleware успел отдать
	// 504 до того, как сервер оборвёт соединение. При слишком маленьком
	// SERVER_WRITE_TIMEOUT запас невозможен — берём WriteTimeout как есть
	// (504 не гарантируется).
	requestTimeout := a.Config.Server.WriteTimeout - 5*time.Second
	if requestTimeout <= 0 {
		requestTimeout = a.Config.Server.WriteTimeout
	}
	r.Use(middleware.Timeout(requestTimeout))

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
		httpSwagger.URL("/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:              ":" + a.Config.Port,
		Handler:           r,
		ReadTimeout:       a.Config.Server.ReadTimeout,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      a.Config.Server.WriteTimeout,
		MaxHeaderBytes:    a.Config.Server.MaxHeaderBytes,
	}

	a.Logger.Info().
		Str("port", a.Config.Port).
		Msg("server initialized")

	return server, nil
}
