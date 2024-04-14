-- name: CreateExercise :one
INSERT INTO exercises (id, name, muscle, sets, repetitions, exercise_duration, instructions, exercise_type, day_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetExercisesByDay :many
SELECT * FROM exercises
WHERE exercises.day_id = $1;