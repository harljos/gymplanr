-- name: CreateDay :one
INSERT INTO days (id, name, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;