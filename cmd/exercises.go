package cmd

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/harljos/gymplanr/internal/database"
)

func (cfg *config) createExercise(name, muscle, instructions string, repetitions int, day database.Day) (database.Exercise, error) {
	exercise, err := cfg.DB.CreateExercise(context.Background(), database.CreateExerciseParams{
		ID:           uuid.New(),
		Name:         name,
		Muscle:       muscle,
		Repetitions:  int32(repetitions),
		Instructions: instructions,
		DayID:        day.ID,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		return database.Exercise{}, err
	}

	return exercise, nil
}
