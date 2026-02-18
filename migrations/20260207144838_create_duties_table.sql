-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS duties
(
    id          UUID    NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    date        VARCHAR NOT NULL,
    in_workdays INT     NOT NULL,
    in_holidays INT     NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS duties;
-- +goose StatementEnd
