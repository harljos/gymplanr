package cmd

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/harljos/gymplanr/internal/database"
)

func (cfg *config) createExercise(name, muscle, instructions, exerciseType, difficulty string, sets, repetitions, duration int, day database.Day) (database.Exercise, error) {
	reps := sql.NullInt32{}
	if repetitions != 0 {
		reps.Int32 = int32(repetitions)
		reps.Valid = true
	}

	setsNum := sql.NullInt32{}
	if sets != 0 {
		setsNum.Int32 = int32(sets)
		setsNum.Valid = true
	}

	durationNum := sql.NullInt32{}
	if duration != 0 {
		durationNum.Int32 = int32(duration)
		durationNum.Valid = true
	}

	exercise, err := cfg.DB.CreateExercise(context.Background(), database.CreateExerciseParams{
		ID:               uuid.New(),
		Name:             name,
		Muscle:           muscle,
		ExerciseType:     exerciseType,
		Difficulty:       difficulty,
		Sets:             setsNum,
		Repetitions:      reps,
		ExerciseDuration: durationNum,
		Instructions:     instructions,
		DayID:            day.ID,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	})
	if err != nil {
		return database.Exercise{}, err
	}

	return exercise, nil
}

func (cfg *config) getExercisesByDay(day database.Day) ([]database.Exercise, error) {
	exercises, err := cfg.DB.GetExercisesByDay(context.Background(), day.ID)
	if err != nil {
		return []database.Exercise{}, err
	}

	return exercises, nil
}

func (cfg *config) updateExercise(muscle, difficulty, exerciseType string, databaseExercise database.Exercise) error {
	exercise, err := cfg.exerciseClient.GetExercise(muscle, difficulty, exerciseType)
	if err != nil {
		return err
	}

	err = cfg.DB.UpdateExercise(context.Background(), database.UpdateExerciseParams{
		Name: exercise.Name,
		Muscle: exercise.Muscle,
		Instructions: exercise.Instructions,
		ExerciseType: exercise.Type,
		Difficulty: exercise.Difficulty,
		UpdatedAt: time.Now().UTC(),
		ID: databaseExercise.DayID,
	})
	if err != nil {
		return err
	}

	return nil
}
