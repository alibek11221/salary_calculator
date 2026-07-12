# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Go HTTP API (chi router) for calculating salary and managing bonuses, salary changes, and duties. PostgreSQL via pgx/v5 with sqlc-generated queries. Swagger UI at `/swagger/`, API mounted at `/api/v1`.

## Commands

A `.env` file is required (copy from `.env.example`); the Makefile loads it. Most make targets have short aliases (shown in `make help`).

```bash
make run-local          # go run ./cmd/main.go
make build              # build to bin/salary_calculator
make fmt                # gofumpt + goimports (-local salary_calculator)
make fmt-check          # formatting check without changes

go test ./...                                        # run all tests (no make target)
go test ./internal/services/calculator/ -run TestName # run a single test

make run-docker         # start Postgres (docker compose)
make run-migrations     # apply goose migrations
make full-reset         # recreate Docker volumes + reapply migrations
make create-migration name=<name>  # new goose migration (auto git-adds it)

make mocks              # regenerate mockgen mocks (go generate ./...)
make sqlc-gen           # regenerate internal/generated/dbstore from queries/
make swagger            # regenerate docs/ from handler annotations
make deps-update        # go mod tidy + verify + vendor
```

Dev tools (goose, sqlc, mockgen, swag, gofumpt, goimports) are installed into `tools/bin/` by `make tools`; targets install them automatically when missing.

Dependencies are **vendored** — after any `go.mod` change run `make deps-update` to refresh `vendor/`.

## Architecture

Wiring flow: `cmd/main.go` → `internal/app.App` (holds Config, DB pool, zerolog Logger, and `Repo` — the sqlc `*dbstore.Queries`) → `internal/http/server` (chi + middleware) → `internal/http/routes`.

**Routes**: each resource (health, salary, bonuses, duties) has a `RouterRegistrar` in `internal/http/routes/` whose `Name()` becomes the URL segment (`/api/v1/bonuses`, etc.). The registrar constructs usecases from `app.Repo` and binds handlers. Adding an endpoint means touching a registrar file (or adding a new one to `NewRoutesRegistrar`).

**Package-per-action layering** — every endpoint gets three small packages:
- `internal/http/handlers/<action>/` — `handler.go` (HTTP + swagger annotations) and `contract.go` declaring the `usecase` interface it consumes
- `internal/usecase/<domain>/<action>/` — `usecase.go` with business logic and `contract.go` declaring the `repo` (and service) interfaces it consumes; domain sentinel errors (e.g. `ErrDuplicateBonus`) live here
- `internal/dto/<action>/` — request/response In/Out types

Interfaces are always defined on the consumer side in `contract.go`, each carrying a `//go:generate mockgen -source=contract.go -destination mocks_test.go` directive; mocks are generated into the package's `_test` package. When you change a `contract.go`, run `make mocks`.

**Database**: schema lives in goose migrations (`migrations/`), queries in `queries/*.sql`; sqlc reads both (see `sqlc.yaml`) and generates `internal/generated/dbstore/` — never edit generated files, change the SQL and run `make sqlc-gen`.

**Salary calculation domain**: production-calendar JSON files (`const/workdays/workdays_YYYY.json`, dir configurable via `WORKDAYS_DIR`) are read by `internal/pkg/http/work_calendar_parser` (LRU-cached). `internal/services/work_days` splits a month's workdays into first half (days 1–15) and second half — salary is paid in two parts. `internal/services/calculator` computes the report; NDFL (income tax) helpers are in `internal/pkg/utils`. `get_salary_report` usecase fetches the latest salary change before the target date and the workday calendar concurrently via errgroup.

**Shared helpers** (`internal/pkg/`): `logging` (zerolog wrapper + chi middleware), `http/response` (JSON response writing), `types.Date`, `pointer`, `utils`.
