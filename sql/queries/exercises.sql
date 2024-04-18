-- name: CreateExercise :one
INSERT INTO exercises (id, name, muscle, sets, repetitions, exercise_duration, instructions, exercise_type, difficulty, day_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetExercisesByDay :many
SELECT * FROM exercises
WHERE exercises.day_id = $1;

-- name: UpdateSets :exec
UPDATE exercises
SET sets = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateReps :exec
UPDATE exercises
SET repetitions = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateDuration :exec
UPDATE exercises
SET exercise_duration = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateExercise :exec
UPDATE exercises
SET name = $1, muscle = $2, instructions = $3, exercise_type = $4, difficulty = $5, updated_at = $6
WHERE id = $7;