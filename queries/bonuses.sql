-- name: ListBonuses :many
SELECT *
FROM bonuses
ORDER BY date;

-- name: GetBonusByDate :one
SELECT *
FROM bonuses
WHERE date = $1;

-- name: InsertBonus :exec
INSERT INTO bonuses (value, date)
VALUES ($1, $2);

-- name: UpdateBonus :exec
UPDATE bonuses
SET value = $2,
    date  = $3
WHERE id = $1;

-- name: DeleteBonus :exec
DELETE
FROM bonuses
WHERE id = $1;
