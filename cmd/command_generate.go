package cmd

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/harljos/gymplanr/internal/database"
)

const (
	difficultyKey   = "difficulty"
	exerciseTypeKey = "exerciseType"
	daysKey         = "days"
	minutesKey      = "minutes"
	cardioKey       = "cardio"
)

func generateCmd(cfg *config, user database.User) error {
	days, err := cfg.getDaysByUser(user)
	if err != nil {
		return err
	}

	if days != nil {
		ok := YesNoPrompt("Continuing will delete your current workout routine. Do you wish to continue", true)
		if ok {
			err = cfg.deleteDays(user)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	results := make(map[string]string)

	exerciseTypePrompt := []string{"strength", "cardio", "both"}

	result, err := SelectPrompt("what would you like your exercise plan to be based around?", exerciseTypePrompt)
	if err != nil {
		return err
	}
	results[exerciseTypeKey] = result

	if results[exerciseTypeKey] == "strength" || results[exerciseTypeKey] == "both" {
		difficultyPrompt := []string{"beginner", "intermediate", "expert"}

		result, err = SelectPrompt("What would you like the difficulty of the strength exercises to be?", difficultyPrompt)
		if err != nil {
			return err
		}
		results[difficultyKey] = result
	}

	daysPrompt := []string{"3", "4", "5", "6"}

	result, err = SelectPrompt("How many days a week do you want to work out?", daysPrompt)
	if err != nil {
		return err
	}
	results[daysKey] = result

	if results[exerciseTypeKey] == "both" {
		minutesPrompt := []string{"30", "45", "60", "75"}

		result, err = SelectPrompt("How many minutes do you want each workout to be (cardio not included)?", minutesPrompt)
		if err != nil {
			return err
		}
		results[minutesKey] = result

		hoursCardioPrompt := []string{"15", "30", "45", "60", "75", "80"}

		result, err = SelectPrompt("How many minutes of cardio do you want to do?", hoursCardioPrompt)
		if err != nil {
			return err
		}
		results[cardioKey] = result
	} else if results[exerciseTypeKey] == "strength" {
		hoursPrompt := []string{"30", "45", "60", "75"}

		result, err = SelectPrompt("How many minutes do you want each workout to be?", hoursPrompt)
		if err != nil {
			return err
		}
		results[minutesKey] = result
	} else {
		hoursCardioPrompt := []string{"15", "30", "45", "60", "75", "80"}

		result, err = SelectPrompt("How many minutes of cardio do you want to do?", hoursCardioPrompt)
		if err != nil {
			return err
		}
		results[cardioKey] = result
	}

	workoutDays, err := getWorkoutDays(results)
	if err != nil {
		return err
	}

	_, err = cfg.createDays(workoutDays, user)
	if err != nil {
		return err
	}

	go generateWorkout(cfg, user, workoutDays, results)

	fmt.Println("Your workout plan has been generated use 'view' command to see it")

	return nil
}

func generateWorkout(cfg *config, user database.User, days []Day, results map[string]string) {
	wg := &sync.WaitGroup{}
	for _, day := range days {
		wg.Add(1)

		go getExercises(cfg, wg, user, day, results)
	}
	wg.Wait()
}

func getExercises(cfg *config, wg *sync.WaitGroup, user database.User, day Day, results map[string]string) {
	defer wg.Done()

	if results[exerciseTypeKey] == "strength" {
		generateStrengthExercises(cfg, user, day, results)
	} else if results[exerciseTypeKey] == "cardio" {
		generateCardioExercise(cfg, user, day, results)
	} else {
		generateStrengthExercises(cfg, user, day, results)
		generateCardioExercise(cfg, user, day, results)
	}
}

func generateStrengthExercises(cfg *config, user database.User, day Day, results map[string]string) {
	for _, muscle := range day.muscles {
		exercise, err := cfg.exerciseClient.GetExercise(muscle, results[difficultyKey], "strength")
		if err != nil {
			log.Printf("couldn't fetch exercise: %v\n", err)
			continue
		}

		databaseDay, err := cfg.getDayByUser(user, day.dayName)
		if err != nil {
			log.Printf("couldn't get day from database: %v\n", err)
			continue
		}

		_, err = cfg.createExercise(exercise.Name, exercise.Muscle, exercise.Instructions, "strength", 3, 10, 0, databaseDay)
		if err != nil {
			log.Printf("couldn't create exercise: %v\n", err)
			continue
		}
	}
}

func generateCardioExercise(cfg *config, user database.User, day Day, results map[string]string) {
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

	minutes, err := strconv.Atoi(results[minutesKey])
	if err != nil {
		log.Printf("couldn'tcovert to int: %v\n", err)
		return
	}

	_, err = cfg.createExercise(exercise.Name, exercise.Muscle, exercise.Instructions, "cardio", 0, 0, minutes, databaseDay)
	if err != nil {
		log.Printf("couldn't create exercise: %v\n", err)
		return
	}
}
