package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gymplanr",
	Short: "Gymplanr is an application designed to create a full workout specifically for you without the hustle and bustle.",
	Long:  "Gymplanr is an application designed to create a full workout specifically for you without the hustle and bustle.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
