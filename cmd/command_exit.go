package cmd

import "os"

func exitCommand(cfg *config) error {
	os.Exit(0)
	return nil
}