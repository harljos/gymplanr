-- name: CreateDay :one
INSERT INTO days (id, name, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetDaysByUser :many
SELECT * FROM days
WHERE days.user_id = $1;

-- name: DeleteDays :exec
DELETE FROM days WHERE days.user_id = $1;