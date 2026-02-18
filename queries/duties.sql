-- name: ListDuties :many
SELECT *
FROM duties
ORDER BY date;

-- name: GetDutyByDate :one
SELECT *
FROM duties
WHERE date = $1;

-- name: InsertDuty :exec
INSERT INTO duties(date, in_workdays, in_holidays)
VALUES ($1, $2, $3);

-- name: UpdateDuty :exec
UPDATE duties
SET date        = $2,
    in_workdays = $3,
    in_holidays = $4
WHERE id = $1;

-- name: DeleteDuty :exec
DELETE
FROM duties
WHERE id = $1;
