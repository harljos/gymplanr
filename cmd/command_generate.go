package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/harljos/gymplanr/internal/database"
)

var (
	confirm         bool
	exerciseType    string
	difficulty      string
	numOfDays       int
	strengthMinutes string
	cardioMinutes   string
)

func generateCmd(cfg *config, user database.User) error {
	days, err := cfg.getDaysByUser(user)
	if err != nil {
		return err
	}

	if days != nil {
		err := huh.NewConfirm().
			Title("Continuing will delete your current workout routine. Do you wish to continue").
			Affirmative("Yes").
			Negative("No").
			Value(&confirm).
			Run()
		if err != nil {
			return err
		}

		if confirm {
			err = cfg.deleteDays(user)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	err = huh.NewSelect[string]().
		Title("what would you like your exercise plan to be based around?").
		Options(
			huh.NewOption("Strength", "strength"),
			huh.NewOption("Cardio", "cardio"),
			huh.NewOption("Both", "both"),
		).
		Value(&exerciseType).
		Run()
	if err != nil {
		return err
	}

	if exerciseType == "strength" || exerciseType == "both" {
		err := huh.NewSelect[string]().
			Title("What would you like the difficulty of the strength exercises to be?").
			Options(
				huh.NewOption("Beginner", "beginner"),
				huh.NewOption("Intermediate", "intermediate"),
				huh.NewOption("Expert", "expert"),
			).
			Value(&difficulty).
			Run()
		if err != nil {
			return err
		}
	}

	err = huh.NewSelect[int]().
		Title("How many days a week do you want to work out?").
		Options(
			huh.NewOption("3", 3),
			huh.NewOption("4", 4),
			huh.NewOption("5", 5),
			huh.NewOption("6", 6),
		).
		Value(&numOfDays).
		Run()
	if err != nil {
		return err
	}

	switch exerciseType {
	case "both":
		result, err := enterIntString("How many minutes do you want each workout to be (cardio not included)?")
		if err != nil {
			return err
		}
		strengthMinutes = result

		result, err = enterIntString("How many minutes of cardio do you want to do?")
		if err != nil {
			return err
		}
		cardioMinutes = result
	case "strength":
		result, err := enterIntString("How many minutes do you want each workout to be?")
		if err != nil {
			return err
		}
		strengthMinutes = result
	default:
		result, err := enterIntString("How many minutes of cardio do you want to do?")
		if err != nil {
			return err
		}
		cardioMinutes = result
	}

	workoutDays, err := getWorkoutDays()
	if err != nil {
		return err
	}

	_, err = cfg.createDays(workoutDays, user)
	if err != nil {
		return err
	}

	go generateWorkout(cfg, user, workoutDays)

	fmt.Println("Your workout plan has been generated use 'view' command to see it")

	return nil
}

func generateWorkout(cfg *config, user database.User, days []Day) {
	wg := &sync.WaitGroup{}
	for _, day := range days {
		wg.Add(1)

		go getExercises(cfg, wg, user, day)
	}
	wg.Wait()
}

func getExercises(cfg *config, wg *sync.WaitGroup, user database.User, day Day) {
	defer wg.Done()

	switch exerciseType {
	case "strength":
		generateStrengthExercises(cfg, user, day)
	case "cardio":
		generateCardioExercise(cfg, user, day)
	default:
		generateStrengthExercises(cfg, user, day)
		generateCardioExercise(cfg, user, day)
	}
}

func generateStrengthExercises(cfg *config, user database.User, day Day) {
	for _, muscle := range day.muscles {
		exercise, err := cfg.exerciseClient.GetExercise(muscle, difficulty, "strength")
		if err != nil {
			log.Printf("couldn't fetch exercise: %v\n", err)
			continue
		}

		databaseDay, err := cfg.getDayByUser(user, day.dayName)
		if err != nil {
			log.Printf("couldn't get day from database: %v\n", err)
			continue
		}

		sets, reps := getSetsAndReps()

		_, err = cfg.createExercise(exercise.Name, exercise.Muscle, exercise.Instructions, "strength", exercise.Difficulty, sets, reps, 0, databaseDay)
		if err != nil {
			log.Printf("couldn't create exercise: %v\n", err)
			continue
		}
	}
}

func generateCardioExercise(cfg *config, user database.User, day Day) {
	exercise, err := cfg.exerciseClient.GetCardioExercise()
	if err != nil {
		log.Printf("couldn't fetch exercise: %v\n", err)
		return
	}

	databaseDay, err := cfg.getDayByUser(user, day.dayName)
	if err != nil {
		log.Printf("couldn't get day from database: %v\n", err)
		return
	}

	minutes, err := strconv.Atoi(cardioMinutes)
	if err != nil {
		log.Printf("couldn't covert to int: %v\n", err)
		return
	}

	_, err = cfg.createExercise(exercise.Name, exercise.Muscle, exercise.Instructions, "cardio", exercise.Difficulty, 0, 0, minutes, databaseDay)
	if err != nil {
		log.Printf("couldn't create exercise: %v\n", err)
		return
	}
}

func getSetsAndReps() (int, int) {
	switch difficulty {
	case "beginner":
		return 3, 6
	case "intermediate":
		return 3, 12
	case "expert":
		return 4, 12
	}

	return 0, 0
}

func enterIntString(s string) (string, error) {
	var value string
	minutesPrompt := huh.NewText().Title(s).Value(&value)

	err := minutesPrompt.Run()
	if err != nil {
		return "", err
	}
	for value == "" {
		fmt.Println("Please enter a numberic value")
		err = minutesPrompt.Run()
		if err != nil {
			return "", err
		}
	}

	_, err = strconv.Atoi(value)
	if err != nil {
		if strings.Contains(err.Error(), "invalid syntax") {
			fmt.Println("Please enter a numeric value")
			return enterIntString(s)
		}
		return "", err
	}

	return value, nil
}
