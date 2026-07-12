-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vacations
(
    id        UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_from DATE NOT NULL,
    date_to   DATE NOT NULL,
    CHECK (date_to >= date_from),
    CONSTRAINT vacations_no_overlap
        EXCLUDE USING gist (daterange(date_from, date_to, '[]') WITH &&)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vacations;
-- +goose StatementEnd
