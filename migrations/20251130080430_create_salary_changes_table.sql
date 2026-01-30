-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS salary_changes
(
    id          UUID             NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    salary      DOUBLE PRECISION NOT NULL CHECK (salary > 0),
    change_from VARCHAR          NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS change_from_uidx ON salary_changes (change_from);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS change_from_uidx;
DROP TABLE IF EXISTS salary_changes;
-- +goose StatementEnd
