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
	if databaseDays == nil {
		fmt.Println("Workout plan has not been found use 'generate' command to get one")
		return nil
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

	viewExercises(cfg, databaseDays[index], user)

	return nil
}

func viewExercises(cfg *config, day database.Day, user database.User) error {
	databaseExercises, err := cfg.getExercisesByDay(day)
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

	index, result, err := SelectPrompt("Select an exercise for further details", exercises)
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
		updatePrompt := []string{"instructions", "change sets", "change reps"}
		switch exercise.Difficulty {
		case "beginner":
			updatePrompt = append(updatePrompt, "harder exercise")
		case "intermediate":
			updatePrompt = append(updatePrompt, "easier exercise", "harder exercise")
		default:
			updatePrompt = append(updatePrompt, "easier exercise")
		}
		updatePrompt = append(updatePrompt, "change exercise", "back", "quit")

		_, result, err = SelectPrompt("Select one", updatePrompt)
		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond)
		switch result {
		case "quit":
			return nil
		case "back":
			return viewExercises(cfg, day, user)
		case "instructions":
			if exercise.Instructions == "" {
				fmt.Println("No instructions found")
				return nil
			}
			fmt.Println(exercise.Instructions)
			return nil
		case "change sets":
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

			fmt.Println("exercise sets updated")
			return viewExercises(cfg, day, user)
		case "change reps":
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
				UpdatedAt:   time.Now().UTC(),
				ID:          exercise.ID,
			})
			if err != nil {
				return err
			}

			fmt.Println("exercise reps updated")
			return viewExercises(cfg, day, user)
		case "change exercise":
			err = cfg.updateExercise(exercise.Muscle, exercise.Difficulty, exercise.ExerciseType, exercise)
			if err != nil {
				return err
			}

			fmt.Println("new exercise generated")
			return viewExercises(cfg, day, user)
		case "easier exercise":
			err = cfg.updateExercise(exercise.Muscle, easierExercise(exercise.Difficulty), exercise.ExerciseType, exercise)
			if err != nil {
				return err
			}

			fmt.Println("easier exercise generated")
			return viewExercises(cfg, day, user)
		case "harder exercise":
			err = cfg.updateExercise(exercise.Muscle, harderExercise(exercise.Difficulty), exercise.ExerciseType, exercise)
			if err != nil {
				return err
			}

			fmt.Println("harder exercise generated")
			return viewExercises(cfg, day, user)
		}
	} else {
		updatePrompt := []string{"instructions", "change cardio time", "change exercise", "back", "quit"}

		_, result, err = SelectPrompt("Select one", updatePrompt)
		if err != nil {
			return err
		}
		
		time.Sleep(time.Millisecond)
		switch result {
		case "quit":
			return nil
		case "back":
			return viewExercises(cfg, day, user)
		case "instructions":
			if exercise.Instructions == "" {
				fmt.Println("No instructions found")
				return nil
			}
			fmt.Println(exercise.Instructions)
			return nil
		case "change cardio time":
			minutes, err := enterInt("Minutes:")
			if err != nil {
				return err
			}

			duration := sql.NullInt32{
				Int32: int32(minutes),
				Valid: true,
			}

			err = cfg.DB.UpdateDuration(context.Background(), database.UpdateDurationParams{
				ExerciseDuration: duration,
				UpdatedAt:        time.Now().UTC(),
				ID:               exercise.ID,
			})
			if err != nil {
				return err
			}

			fmt.Println("exercise time updated")
			return viewExercises(cfg, day, user)
		case "change exercise":
			err = cfg.updateExercise(exercise.Muscle, exercise.Difficulty, exercise.ExerciseType, exercise)
			if err != nil {
				return err
			}

			fmt.Println("new exercise generated")
			return viewExercises(cfg, day, user)
		}
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

func easierExercise(difficulty string) string {
	if difficulty == "intermediate" {
		return "beginner"
	}
	return "intermediate"
}

func harderExercise(difficulty string) string {
	if difficulty == "beinner" {
		return "intermediate"
	}
	return "expert"
}
