package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
	time.Sleep(time.Millisecond)
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
			exercises = append(exercises, fmt.Sprintf("%s %v sets, %v reps", exercise.Name, exercise.Sets.Int32, exercise.Repetitions.Int32))
		} else {
			exercises = append(exercises, fmt.Sprintf("%s %v minutes", exercise.Name, exercise.ExerciseDuration.Int32))
		}
	}
	exercises = append(exercises, "back", "quit")

	index, result, err = SelectPrompt("Select an exercise for further details", exercises)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond)
	if result == "quit" {
		return nil
	}
	if result == "back" {
		return viewCmd(cfg, user)
	}

	exercise := databaseExercises[index]

	if exercise.ExerciseType == "strength" {
		updatePrompt := []string{"instructions", "change sets", "change reps", "quit"}

		_, result, err = SelectPrompt("What would you like to update?", updatePrompt)
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond)
		if result == "quit" {
			return nil
		}
		if result == "instructions" {
			if exercise.Instructions == "" {
				fmt.Println("No instructions found")
				return nil
			}
			fmt.Println(exercise.Instructions)
			return nil
		}
		if result == "change sets" {
			sets, err := enterInt("Sets:")
			if err != nil {
				return err
			}

			setsNum := sql.NullInt32{
				Int32: int32(sets),
				Valid: true,
			}

			err = cfg.DB.UpdateSets(context.Background(), database.UpdateSetsParams{
				Sets:      setsNum,
				UpdatedAt: time.Now().UTC(),
				ID:        exercise.ID,
			})
			if err != nil {
				return err
			}
		}
		if result == "change reps" {
			reps, err := enterInt("Reps:")
			if err != nil {
				return err
			}

			repsNum := sql.NullInt32{
				Int32: int32(reps),
				Valid: true,
			}

			err = cfg.DB.UpdateReps(context.Background(), database.UpdateRepsParams{
				Repetitions: repsNum,
				UpdatedAt: time.Now().UTC(),
				ID: exercise.ID,
			})
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("minutes: %v\ninstructions: %s\n", exercise.ExerciseDuration.Int32, exercise.Instructions)
	}

	return nil
}

func enterInt(s string) (int, error) {
	stringNum := StringPrompt(s)
	for stringNum == "" {
		fmt.Println("Please enter a numberic value")
		stringNum = StringPrompt(s)
	}

	num, err := strconv.Atoi(stringNum)
	if err != nil {
		if strings.Contains(err.Error(), "invalid syntax") {
			fmt.Println("Please enter a numeric value")
			return enterInt(s)
		}
		log.Printf("couldn't covert to int: %v\n", err)
		return 0, nil
	}

	return num, nil
}
