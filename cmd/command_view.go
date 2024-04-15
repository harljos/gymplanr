package cmd

import (
	"fmt"

	"github.com/harljos/gymplanr/internal/database"
)

func viewCmd(cfg *config, user database.User) error {
	databaseDays, err := cfg.getDaysByUser(user)
	if err != nil {
		return err
	}

	days := []string{}
	for _, day := range databaseDays {
		days = append(days, day.Name)
	}
	days = append(days, "quit")

	index, result, err := SelectPrompt("Select a day you would like to view", days)
	if err != nil {
		return err
	}
	if result == "quit" {
		return nil
	}

	databaseExercises, err := cfg.getExercisesByDay(databaseDays[index])
	if err != nil {
		return err
	}

	exercises := []string{}
	for _, exercise := range databaseExercises {
		if exercise.ExerciseType == "strength" {
			exercises = append(exercises, fmt.Sprintf("%s %v x %v", exercise.Name, exercise.Sets.Int32, exercise.Repetitions.Int32))
		} else {
			exercises = append(exercises, fmt.Sprintf("%s %v minutes", exercise.Name, exercise.ExerciseDuration.Int32))
		}
	}
	exercises = append(exercises, "back", "quit")

	index, result, err = SelectPrompt("Select an exercise for further details", exercises)
	if err != nil {
		return err
	}
	if result == "quit" {
		return nil
	}
	if result == "back" {
		return viewCmd(cfg, user)
	}

	exercise := databaseExercises[index]

	fmt.Printf("sets: %v\nreps: %v\ninstructions: %s\n", exercise.Sets.Int32, exercise.Repetitions.Int32, exercise.Instructions)

	return nil
}
