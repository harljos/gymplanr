package cmd

import "github.com/harljos/gymplanr/internal/database"

func generateCmd(cfg *config, user database.User) error {
	results := make(map[string]string)

	difficulty := []string{"beginner", "intermediate", "expert"}

	result, err := SelectPrompt("What would you like the difficulty of the exercises to be?", difficulty)
	if err != nil {
		return err
	}
	results["difficulty"] = result

	exerciseType := []string{"strength", "cardio", "both"}

	result, err = SelectPrompt("what would you like your exercise plan to be based around?", exerciseType)
	if err != nil {
		return err
	}
	results["exerciseType"] = result

	days := []string{"3", "4", "5", "6"}

	result, err = SelectPrompt("How many days a week do you want to work out?", days)
	if err != nil {
		return err
	}
	results["days"] = result

	hours := []string{"30", "45", "60", "75"}

	result, err = SelectPrompt("How many minutes do you want each workout to be?", hours)
	if err != nil {
		return err
	}
	results["hours"] = result

	workoutDays := getWorkoutDays(results)

	_, err = cfg.createDays(workoutDays, user)
	if err != nil {
		return err
	}

	return nil
}

func getWorkoutDays(results map[string]string) []string {
	days, ok := results["days"]
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
