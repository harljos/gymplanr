package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

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

	_, result, err := SelectPrompt("what would you like your exercise plan to be based around?", exerciseTypePrompt)
	if err != nil {
		return err
	}
	results[exerciseTypeKey] = result

	if results[exerciseTypeKey] == "strength" || results[exerciseTypeKey] == "both" {
		difficultyPrompt := []string{"beginner", "intermediate", "expert"}

		_, result, err = SelectPrompt("What would you like the difficulty of the strength exercises to be?", difficultyPrompt)
		if err != nil {
			return err
		}
		results[difficultyKey] = result
	}

	daysPrompt := []string{"3", "4", "5", "6"}

	_, result, err = SelectPrompt("How many days a week do you want to work out?", daysPrompt)
	if err != nil {
		return err
	}
	results[daysKey] = result

	time.Sleep(time.Millisecond)
	switch results[exerciseTypeKey] {
	case "both":
		result, err := enterIntString("How many minutes do you want each workout to be (cardio not included)?")
		if err != nil {
			return err
		}
		results[minutesKey] = result

		result, err = enterIntString("How many minutes of cardio do you want to do?")
		if err != nil {
			return err
		}
		results[cardioKey] = result
	case "strength":
		result, err = enterIntString("How many minutes do you want each workout to be?")
		if err != nil {
			return err
		}
		results[minutesKey] = result
	default:
		result, err = enterIntString("How many minutes of cardio do you want to do?")
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

	switch results[exerciseTypeKey] {
	case "strength":
		generateStrengthExercises(cfg, user, day, results)
	case "cardio":
		generateCardioExercise(cfg, user, day, results)
	default:
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

		sets, reps := getSetsAndReps(results)

		_, err = cfg.createExercise(exercise.Name, exercise.Muscle, exercise.Instructions, "strength", sets, reps, 0, databaseDay)
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

	minutes, err := strconv.Atoi(results[cardioKey])
	if err != nil {
		log.Printf("couldn't covert to int: %v\n", err)
		return
	}

	_, err = cfg.createExercise(exercise.Name, exercise.Muscle, exercise.Instructions, "cardio", 0, 0, minutes, databaseDay)
	if err != nil {
		log.Printf("couldn't create exercise: %v\n", err)
		return
	}
}

func getSetsAndReps(results map[string]string) (int, int) {
	switch results[difficultyKey] {
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
	stringNum := StringPrompt(s)
	for stringNum == "" {
		fmt.Println("Please enter a numberic value")
		stringNum = StringPrompt(s)
	}

	_, err := strconv.Atoi(stringNum)
	if err != nil {
		if strings.Contains(err.Error(), "invalid syntax") {
			fmt.Println("Please enter a numeric value")
			return enterIntString(s)
		}
		log.Printf("couldn't covert to int: %v\n", err)
		return "", nil
	}

	return stringNum, nil
}
