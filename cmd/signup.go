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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
