package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var signUpCmd = &cobra.Command{
	Use:   "signUp",
	Short: "creates an account for gymplanr",
	Long:  "creates an account for gymplanr",
	Run: func(cmd *cobra.Command, args []string) {
		user := StringPrompt("Username:")
		fmt.Printf("Hello %s\n", user)

		password := PasswordPrompt("Password:")
		fmt.Printf("This is your password %s\n", password)
	},
}

func init() {
	rootCmd.AddCommand(signUpCmd)
}
