package cmd

import "fmt"

func generateCmd(cfg *config) error {
	results := []string{}

	difficulty := []string{"beginner", "intermediate", "expert"}

	result, err := SelectPrompt("What would you like the difficulty of the excercises to be?", difficulty)
	if err != nil {
		return err
	}
	results = append(results, result)

	excerciseType := []string{"strength", "cardio", "both"}

	result, err = SelectPrompt("what would you like your excercise plan to be based around?", excerciseType)
	if err != nil {
		return err
	}
	results = append(results, result)

	days := []string{"3", "4", "5", "6"}

	result, err = SelectPrompt("How many days a week do you want to work out?", days)
	if err != nil {
		return err
	}
	results = append(results, result)

	hours := []string{"30", "45", "60", "75"}

	result, err = SelectPrompt("How many minutes do you want each workout to be?", hours)
	if err != nil {
		return err
	}
	results = append(results, result)

	muscle := []string{"abdominals", "abductors", "adductors", "biceps", "calves",
		"chest", "forearms", "glutes", "hamstrings", "lats", "lower back",
		"middle back", "quadriceps", "traps", "triceps", "none"}

	result, err = SelectPrompt("Pick a muscle to focus on", muscle)
	if err != nil {
		return err
	}
	results = append(results, result)

	fmt.Println(results)

	return nil
}
