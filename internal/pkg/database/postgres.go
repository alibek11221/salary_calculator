package database

import (
	"context"
	"fmt"
	"salary_calculator/internal/config"
	"salary_calculator/internal/pkg/logging"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
	logger logging.Logger
}

func NewPostgresConnection(cfg *config.Config, logger logging.Logger) (*DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	ctx := context.Background()
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.Database.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.Database.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.Database.ConnMaxLifetime / 2

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().
		Str("host", cfg.Database.Host).
		Str("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Int("max_open_conns", cfg.Database.MaxOpenConns).
		Int("max_idle_conns", cfg.Database.MaxIdleConns).
		Dur("conn_max_lifetime", cfg.Database.ConnMaxLifetime).
		Msg("database connection established")

	return &DB{Pool: pool, logger: logger}, nil
}

func (db *DB) Close() {
	if db == nil || db.Pool == nil {
		return
	}
	db.logger.Info().Msg("closing database connection")
	db.Pool.Close()
}

func (db *DB) HealthCheck(ctx context.Context) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return db.Ping(ctx)
}

func (db *DB) GetPool() *pgxpool.Pool {
	if db == nil {
		return nil
	}
	return db.Pool
}
