package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"salary_calculator/internal/app"
	"salary_calculator/internal/config"
	"salary_calculator/internal/http/server"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

// @title Salary Calculator API
// @version 1.0
// @description API для расчета зарплаты, управления бонусами и изменениями оклада.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.GetConfig()

	a, err := app.New(cfg)
	if err != nil {
		panic("cannot initialize app: " + err.Error())
	}

	logger := a.Logger

	srv, err := server.NewServer(a)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server")
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("server shutdown failed")
		}
	}()

	logger.Info().
		Str("port", cfg.Port).
		Msg("starting server")

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Error().Err(err).Msg("server stopped unexpectedly")
	}
}
