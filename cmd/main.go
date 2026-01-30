package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"salary_calculator/internal/app"
	"salary_calculator/internal/config"
	"salary_calculator/internal/http/server"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	env := os.Getenv("APP_ENV")
	if env == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
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
		log.Fatal().Err(err).Msg("cannot initialize app")
		os.Exit(1)
	}

	srv, err := server.NewServer(a)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server shutdown failed")
		}
	}()

	log.Info().
		Str("port", cfg.Port).
		Msg("starting server")

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("server stopped unexpectedly")
	}
}
