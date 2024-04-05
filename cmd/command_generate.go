package cmd

import "fmt"

func generateCmd(cfg *config) error {
	muscles := []string{"bicep", "chest", "tricep"}

	result, err := SelectPrompt("What muscle would you like to focus on", muscles)
	if err != nil {
		return err
	}

	fmt.Printf("You chose %s\n", result)
	return nil
}
