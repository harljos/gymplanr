package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/harljos/gymplanr/internal/database"
)

const (
	difficultyKey   = "difficulty"
	exerciseTypeKey = "exerciseType"
	daysKey         = "days"
	hoursKey        = "hours"
)

func generateCmd(cfg *config, user database.User) error {
	days, err := cfg.getDaysByUser(user)
	if err != nil {
		return err
	}

	if days != nil {
		err = cfg.deleteDays(user)
		if err != nil {
			return err
		}
	}

	results := make(map[string]string)

	difficultyPrompt := []string{"beginner", "intermediate", "expert"}

	result, err := SelectPrompt("What would you like the difficulty of the exercises to be?", difficultyPrompt)
	if err != nil {
		return err
	}
	results[difficultyKey] = result

	exerciseTypePrompt := []string{"strength", "cardio", "both"}

	result, err = SelectPrompt("what would you like your exercise plan to be based around?", exerciseTypePrompt)
	if err != nil {
		return err
	}
	results[exerciseTypeKey] = result

	daysPrompt := []string{"3", "4", "5", "6"}

	result, err = SelectPrompt("How many days a week do you want to work out?", daysPrompt)
	if err != nil {
		return err
	}
	results[daysKey] = result

	hoursPrompt := []string{"30", "45", "60", "75"}

	result, err = SelectPrompt("How many minutes do you want each workout to be?", hoursPrompt)
	if err != nil {
		return err
	}
	results[hoursKey] = result

	workoutDays := getWorkoutDays(results)

	databaseDays, err := cfg.createDays(workoutDays, user)
	if err != nil {
		return err
	}

	go generateWorkout(cfg, databaseDays, results)

	fmt.Println("Your workout plan has been generated")

	return nil
}

func generateWorkout(cfg *config, days []database.Day, results map[string]string) {
	wg := &sync.WaitGroup{}
	for _, day := range days {
		wg.Add(1)

		go getExercises(cfg, wg, day, results)
	}
	wg.Wait()
}

func getExercises(cfg *config, wg *sync.WaitGroup, day database.Day, results map[string]string) {
	defer wg.Done()

	muscles := []string{"chest", "shoulders", "middle_back", "glutes", "hamstrings", "quadriceps"}

	for _, muscle := range muscles {
		exercise, err := cfg.exerciseClient.GetExercise(muscle, results[difficultyKey], results[exerciseTypeKey])
		if err != nil {
			log.Printf("couldn't fetch exercise: %v\n", err)
			continue
		}

		_, err = cfg.createExercise(exercise.Name, exercise.Muscle, exercise.Instructions, 10, day)
		if err != nil {
			log.Printf("couldn't create exercise: %v\n", err)
			continue
		}
	}
}

func getWorkoutDays(results map[string]string) []string {
	days, ok := results[daysKey]
	if !ok {
		return []string{}
	}

	if days == "3" {
		return []string{"Monday", "Wednesday", "Friday"}
	}
	if days == "4" {
		return []string{"Monday", "Tuesday", "Thursday", "Friday"}
	}
	if days == "5" {
		return []string{"Monday", "Tuesday", "Wednesday", "Friday", "Saturday"}
	}
	if days == "6" {
		return []string{"Monday", "Tuesday", "Wednesday", "Thurday", "Friday", "Saturday"}
	}
	return []string{}
}
