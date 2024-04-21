package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/harljos/gymplanr/internal/database"
)

var (
	databaseDay      database.Day
	databaseExercise database.Exercise
	updateOption     string
)

func viewCmd(cfg *config, user database.User) error {
	databaseDays, err := cfg.getDaysByUser(user)
	if err != nil {
		return err
	}
	if databaseDays == nil {
		fmt.Println("Workout plan has not been found use the 'generate' command to get one")
		return nil
	}

	days := []huh.Option[database.Day]{}
	daysPrompt := huh.NewSelect[database.Day]().
		Title("Select a day you would like to view").
		Value(&databaseDay)
	for _, day := range databaseDays {
		days = append(days, huh.NewOption(day.Name, day))
	}
	days = append(days, huh.NewOption("Quit", database.Day{Name: "quit"}))

	err = daysPrompt.Options(days...).Run()
	if err != nil {
		return err
	}
	if databaseDay.Name == "quit" {
		return nil
	}

	viewExercises(cfg, databaseDay, user)

	return nil
}

func viewExercises(cfg *config, day database.Day, user database.User) error {
	databaseExercises, err := cfg.getExercisesByDay(day)
	if err != nil {
		return err
	}

	exercises := []huh.Option[database.Exercise]{}
	exercisesPrompt := huh.NewSelect[database.Exercise]().
		Title("Select an exercise for further details").
		Value(&databaseExercise)
	for _, exercise := range databaseExercises {
		if exercise.ExerciseType == "strength" {
			exercises = append(exercises, huh.NewOption(fmt.Sprintf("%s %v sets, %v reps", exercise.Name, exercise.Sets.Int32, exercise.Repetitions.Int32), exercise))
		} else {
			exercises = append(exercises, huh.NewOption(fmt.Sprintf("%s %v minutes", exercise.Name, exercise.ExerciseDuration.Int32), exercise))
		}
	}
	exercises = append(exercises, huh.NewOption("Back", database.Exercise{Name: "back"}), huh.NewOption("Quit", database.Exercise{Name: "quit"}))

	err = exercisesPrompt.Options(exercises...).Run()
	if err != nil {
		return err
	}
	if databaseExercise.Name == "quit" {
		return nil
	}
	if databaseExercise.Name == "back" {
		return viewCmd(cfg, user)
	}

	if databaseExercise.ExerciseType == "strength" {
		updateOptions := []huh.Option[string]{huh.NewOption("Instructions", "instructions"), huh.NewOption("Change Sets", "change sets"), huh.NewOption("Change Reps", "change reps")}
		updatePrompt := huh.NewSelect[string]().
			Title("Select one").
			Value(&updateOption)
		switch databaseExercise.Difficulty {
		case "beginner":
			updateOptions = append(updateOptions, huh.NewOption("Harder Exercise", "harder exercise"))
		case "intermediate":
			updateOptions = append(updateOptions, huh.NewOption("Easier Exercise", "easier exercise"), huh.NewOption("Harder Exercise", "harder exercise"))
		default:
			updateOptions = append(updateOptions, huh.NewOption("Easier Exercise", "easier exercise"))
		}
		updateOptions = append(updateOptions, huh.NewOption("Change Exercise", "change exercise"), huh.NewOption("Back", "back"), huh.NewOption("Quit", "quit"))

		err = updatePrompt.Options(updateOptions...).Run()
		if err != nil {
			return err
		}

		switch updateOption {
		case "quit":
			return nil
		case "back":
			return viewExercises(cfg, day, user)
		case "instructions":
			if databaseExercise.Instructions == "" {
				fmt.Println("No instructions found")
				return nil
			}
			fmt.Println(databaseExercise.Instructions)
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
				ID:        databaseExercise.ID,
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
				ID:          databaseExercise.ID,
			})
			if err != nil {
				return err
			}

			fmt.Println("exercise reps updated")
			return viewExercises(cfg, day, user)
		case "change exercise":
			err = cfg.updateExercise(databaseExercise.Muscle, databaseExercise.Difficulty, databaseExercise.ExerciseType, databaseExercise)
			if err != nil {
				return err
			}

			fmt.Println("new exercise generated")
			return viewExercises(cfg, day, user)
		case "easier exercise":
			err = cfg.updateExercise(databaseExercise.Muscle, easierExercise(databaseExercise.Difficulty), databaseExercise.ExerciseType, databaseExercise)
			if err != nil {
				return err
			}

			fmt.Println("easier exercise generated")
			return viewExercises(cfg, day, user)
		case "harder exercise":
			err = cfg.updateExercise(databaseExercise.Muscle, harderExercise(databaseExercise.Difficulty), databaseExercise.ExerciseType, databaseExercise)
			if err != nil {
				return err
			}

			fmt.Println("harder exercise generated")
			return viewExercises(cfg, day, user)
		}
	} else {
		updatePrompt := []string{"instructions", "change cardio time", "change exercise", "back", "quit"}

		_, result, err := SelectPrompt("Select one", updatePrompt)
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
			if databaseExercise.Instructions == "" {
				fmt.Println("No instructions found")
				return nil
			}
			fmt.Println(databaseExercise.Instructions)
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
				ID:               databaseExercise.ID,
			})
			if err != nil {
				return err
			}

			fmt.Println("exercise time updated")
			return viewExercises(cfg, day, user)
		case "change exercise":
			err = cfg.updateExercise(databaseExercise.Muscle, databaseExercise.Difficulty, databaseExercise.ExerciseType, databaseExercise)
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
