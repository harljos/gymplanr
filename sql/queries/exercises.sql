-- name: CreateExercise :one
INSERT INTO exercises (id, name, muscle, repetitions, instructions, day_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;