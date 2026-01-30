-- name: ListChanges :many
SELECT *
FROM salary_changes
ORDER BY id;

-- name: GetChangeByDate :one
SELECT id
FROM salary_changes
WHERE change_from = $1
ORDER BY id;


-- name: InsertChanges :copyfrom
INSERT INTO salary_changes (salary, change_from)
VALUES ($1, $2);

-- name: InsertChange :exec
INSERT INTO salary_changes (salary, change_from)
VALUES ($1, $2);

-- name: UpdateChange :exec
UPDATE salary_changes SET
salary = $2,
change_from = $3
WHERE id = $1;

-- name: DeleteChange :exec
DELETE FROM salary_changes WHERE id = $1;
