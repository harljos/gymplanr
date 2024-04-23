package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/harljos/gymplanr/internal/database"
)

var (
	databaseDay      database.Day
	databaseExercise database.Exercise
	updateOption     string
	stringInt        string
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
	if len(databaseDays) == 1 {
		return viewExercises(cfg, databaseDays[0], user, len(databaseDays))
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

	viewExercises(cfg, databaseDay, user, len(databaseDays))

	return nil
}

func viewExercises(cfg *config, day database.Day, user database.User, numOfDays int) error {
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
	if numOfDays != 1 {
		exercises  = append(exercises, huh.NewOption("Back", database.Exercise{Name: "back"}))
	}
	exercises = append(exercises, huh.NewOption("Quit", database.Exercise{Name: "quit"}))

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

	intPrompt := huh.NewInput().
		Title("How many mintues do you want to do strength exercises for?").
		Validate(isInt).
		Value(&stringInt)

	if databaseExercise.ExerciseType == "strength" {
		updateOptions := []huh.Option[string]{
			huh.NewOption("Instructions", "instructions"),
			huh.NewOption("Change Sets", "change sets"),
			huh.NewOption("Change Reps", "change reps"),
		}

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
			return viewExercises(cfg, day, user, numOfDays)
		case "instructions":
			if databaseExercise.Instructions == "" {
				fmt.Println("No instructions found")
				return nil
			}
			fmt.Println(databaseExercise.Instructions)
			return nil
		case "change sets":
			err = intPrompt.Run()
			if err != nil {
				return err
			}

			sets, err := strconv.Atoi(stringInt)
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
			return viewExercises(cfg, day, user, numOfDays)
		case "change reps":
			err = intPrompt.Run()
			if err != nil {
				return err
			}

			reps, err := strconv.Atoi(stringInt)
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
			return viewExercises(cfg, day, user, numOfDays)
		case "change exercise":
			err = cfg.updateExercise(databaseExercise.Muscle, databaseExercise.Difficulty, databaseExercise.ExerciseType, databaseExercise)
			if err != nil {
				return err
			}

			fmt.Println("new exercise generated")
			return viewExercises(cfg, day, user, numOfDays)
		case "easier exercise":
			err = cfg.updateExercise(databaseExercise.Muscle, easierExercise(databaseExercise.Difficulty), databaseExercise.ExerciseType, databaseExercise)
			if err != nil {
				return err
			}

			fmt.Println("easier exercise generated")
			return viewExercises(cfg, day, user, numOfDays)
		case "harder exercise":
			err = cfg.updateExercise(databaseExercise.Muscle, harderExercise(databaseExercise.Difficulty), databaseExercise.ExerciseType, databaseExercise)
			if err != nil {
				return err
			}

			fmt.Println("harder exercise generated")
			return viewExercises(cfg, day, user, numOfDays)
		}
	} else {
		updateOptions := []huh.Option[string]{
			huh.NewOption("Instructions", "instructions"),
			huh.NewOption("Change cardio time", "change cardio time"),
			huh.NewOption("Change exercise", "change exercise"),
			huh.NewOption("Back", "back"),
			huh.NewOption("Quit", "quit"),
		}

		err = huh.NewSelect[string]().
			Title("Select one").
			Value(&updateOption).
			Options(updateOptions...).
			Run()
		if err != nil {
			return err
		}

		switch updateOption {
		case "quit":
			return nil
		case "back":
			return viewExercises(cfg, day, user, numOfDays)
		case "instructions":
			if databaseExercise.Instructions == "" {
				fmt.Println("No instructions found")
				return nil
			}
			fmt.Println(databaseExercise.Instructions)
			return nil
		case "change cardio time":
			err = intPrompt.Run()
			if err != nil {
				return err
			}

			minutes, err := strconv.Atoi(stringInt)
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
			return viewExercises(cfg, day, user, numOfDays)
		case "change exercise":
			err = cfg.updateExercise(databaseExercise.Muscle, databaseExercise.Difficulty, databaseExercise.ExerciseType, databaseExercise)
			if err != nil {
				return err
			}

			fmt.Println("new exercise generated")
			return viewExercises(cfg, day, user, numOfDays)
		}
	}

	return nil
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
