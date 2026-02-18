-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bonuses
(
    id          UUID             NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    value       DOUBLE PRECISION NOT NULL CHECK (value > 0),
    date        VARCHAR          NOT NULL
);
CREATE UNIQUE INDEX bonuses_date_idx ON bonuses (date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX bonuses_date_idx;
DROP TABLE IF EXISTS bonuses;
-- +goose StatementEnd
