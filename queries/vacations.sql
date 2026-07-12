-- name: ListVacations :many
SELECT *
FROM vacations
ORDER BY date_from;

-- name: GetVacationsInRange :many
-- daterange && daterange задействует gist-индекс vacations_no_overlap.
SELECT *
FROM vacations
WHERE daterange(date_from, date_to, '[]') && daterange(sqlc.arg(range_start), sqlc.arg(range_end), '[]')
ORDER BY date_from;

-- name: InsertVacation :exec
INSERT INTO vacations (date_from, date_to)
VALUES ($1, $2);

-- name: UpdateVacation :exec
UPDATE vacations
SET date_from = $2,
    date_to   = $3
WHERE id = $1;

-- name: DeleteVacation :exec
DELETE
FROM vacations
WHERE id = $1;
