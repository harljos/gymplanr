package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var signUpCmd = &cobra.Command{
	Use:   "signUp",
	Short: "creates an account for gymplanr",
	Long:  "creates an account for gymplanr",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := connectToDB()
		if err != nil {
			log.Fatal(err)
		}

		username := usernamePrompt()
		password := passwordPrompt()

		user, err := cfg.createUserHandler(username, password)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Thank you %s for signing up for gymplanr!\n", user.Username.String)
	},
}

func init() {
	rootCmd.AddCommand(signUpCmd)
}
