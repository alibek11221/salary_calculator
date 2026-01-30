package app

import (
	"salary_calculator/internal/config"
	"salary_calculator/internal/pkg/database"
)

type App struct {
	Config *config.Config
	DB     *database.DB
}

func New(cfg *config.Config) (*App, error) {
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config: cfg,
		DB:     db,
	}, nil
}
