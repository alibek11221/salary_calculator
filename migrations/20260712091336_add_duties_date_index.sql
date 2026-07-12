-- +goose Up
-- +goose StatementBegin
-- Дедупликация перед UNIQUE-индексом: на каждую дату остаётся одна строка
-- (детерминированно — с минимальным id).
DELETE
FROM duties a
    USING duties b
WHERE a.date = b.date
  AND a.id > b.id;
-- +goose StatementEnd
-- +goose StatementBegin
-- UNIQUE намеренно: код всегда подразумевал одно дежурство на месяц
-- (GetDutyByDate :one, sentinel ErrDuplicateDuty). Индекс фиксирует этот
-- инвариант на уровне БД и оживляет ветку 409 в POST /duties.
CREATE UNIQUE INDEX IF NOT EXISTS duties_date_idx ON duties (date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS duties_date_idx;
-- +goose StatementEnd
