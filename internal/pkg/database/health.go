package database

import (
	"context"
)

type HealthChecker struct {
	db *DB
}

func NewHealthChecker(db *DB) *HealthChecker {
	return &HealthChecker{db: db}
}

func (hc *HealthChecker) HealthCheck(ctx context.Context) error {
	return hc.db.HealthCheck(ctx)
}

func (hc *HealthChecker) Name() string {
	return "database"
}
