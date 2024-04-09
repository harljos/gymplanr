package cmd

import (
	"os"

	"github.com/harljos/gymplanr/internal/database"
)

func exitCmd(cfg *config, user database.User) error {
	os.Exit(0)
	return nil
}
