-- +goose Up
-- +goose StatementBegin
INSERT INTO salary_changes (salary, change_from)
VALUES (252000, '2024_04'),
       (328000, '2025_03'),
       (349000, '2025_09');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM salary_changes
where change_from IN ('04_2024', '03_2025', '09_2025')
-- +goose StatementEnd
