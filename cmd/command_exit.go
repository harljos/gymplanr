package cmd

import "os"

func exitCmd(cfg *config) error {
	os.Exit(0)
	return nil
}