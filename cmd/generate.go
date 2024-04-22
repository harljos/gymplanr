package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generates workout plan",
	Long:  "Generates workout plan",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := connectToDB()
		if err != nil {
			log.Fatal(err)
		}

		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}

		user, err := cfg.getLocalUser(hostname)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				user, err = cfg.createLocalUser(hostname)
				if err != nil {
					log.Fatal(err)
				}
				generateCmd(&cfg, user)
			} else {
				log.Fatal(err)
			}
		} else {
			generateCmd(&cfg, user)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCommand)
}
