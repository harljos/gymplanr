// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: exercises.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createExercise = `-- name: CreateExercise :one
INSERT INTO exercises (id, name, muscle, sets, repetitions, exercise_duration, instructions, exercise_type, day_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, name, muscle, sets, repetitions, exercise_duration, instructions, exercise_type, day_id, created_at, updated_at
`

type CreateExerciseParams struct {
	ID               uuid.UUID
	Name             string
	Muscle           string
	Sets             sql.NullInt32
	Repetitions      sql.NullInt32
	ExerciseDuration sql.NullInt32
	Instructions     string
	ExerciseType     string
	DayID            uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) CreateExercise(ctx context.Context, arg CreateExerciseParams) (Exercise, error) {
	row := q.db.QueryRowContext(ctx, createExercise,
		arg.ID,
		arg.Name,
		arg.Muscle,
		arg.Sets,
		arg.Repetitions,
		arg.ExerciseDuration,
		arg.Instructions,
		arg.ExerciseType,
		arg.DayID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Exercise
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Muscle,
		&i.Sets,
		&i.Repetitions,
		&i.ExerciseDuration,
		&i.Instructions,
		&i.ExerciseType,
		&i.DayID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getExercisesByDay = `-- name: GetExercisesByDay :many
SELECT id, name, muscle, sets, repetitions, exercise_duration, instructions, exercise_type, day_id, created_at, updated_at FROM exercises
WHERE exercises.day_id = $1
`

func (q *Queries) GetExercisesByDay(ctx context.Context, dayID uuid.UUID) ([]Exercise, error) {
	rows, err := q.db.QueryContext(ctx, getExercisesByDay, dayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Exercise
	for rows.Next() {
		var i Exercise
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Muscle,
			&i.Sets,
			&i.Repetitions,
			&i.ExerciseDuration,
			&i.Instructions,
			&i.ExerciseType,
			&i.DayID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateReps = `-- name: UpdateReps :exec
UPDATE exercises
SET repetitions = $1, updated_at = $2
WHERE id = $3
`

type UpdateRepsParams struct {
	Repetitions sql.NullInt32
	UpdatedAt   time.Time
	ID          uuid.UUID
}

func (q *Queries) UpdateReps(ctx context.Context, arg UpdateRepsParams) error {
	_, err := q.db.ExecContext(ctx, updateReps, arg.Repetitions, arg.UpdatedAt, arg.ID)
	return err
}

const updateSets = `-- name: UpdateSets :exec
UPDATE exercises
SET sets = $1, updated_at = $2
WHERE id = $3
`

type UpdateSetsParams struct {
	Sets      sql.NullInt32
	UpdatedAt time.Time
	ID        uuid.UUID
}

func (q *Queries) UpdateSets(ctx context.Context, arg UpdateSetsParams) error {
	_, err := q.db.ExecContext(ctx, updateSets, arg.Sets, arg.UpdatedAt, arg.ID)
	return err
}
