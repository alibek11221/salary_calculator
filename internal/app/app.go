package app

import (
	"salary_calculator/internal/config"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/pkg/database"
	"salary_calculator/internal/pkg/logging"
)

type App struct {
	Config *config.Config
	DB     *database.DB
	Logger logging.Logger
	Repo   *dbstore.Queries
}

func New(cfg *config.Config) (*App, error) {
	logger := logging.New(cfg.Env == "production")
	db, err := database.NewPostgresConnection(cfg, logger)
	if err != nil {
		return nil, err
	}

	repo := dbstore.New(db)

	return &App{
		Config: cfg,
		DB:     db,
		Logger: logger,
		Repo:   repo,
	}, nil
}
