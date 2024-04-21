package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "logsin user",
	Long:  "logsin user",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := connectToDB()
		if err != nil {
			log.Fatal(err)
		}

		username := usernamePrompt()
		password := passwordPrompt()

		user, err := cfg.loginUserHandler(username, password)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Welcome to Gymplanr %s!\n", user.Username.String)
		startRepl(&cfg, user)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
