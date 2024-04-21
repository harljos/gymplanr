// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Day struct {
	ID        uuid.UUID
	Name      string
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Exercise struct {
	ID               uuid.UUID
	Name             string
	Muscle           string
	Sets             sql.NullInt32
	Repetitions      sql.NullInt32
	ExerciseDuration sql.NullInt32
	Instructions     string
	ExerciseType     string
	Difficulty       string
	DayID            uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type User struct {
	ID        uuid.UUID
	Username  sql.NullString
	Password  sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	Hostname  sql.NullString
}
