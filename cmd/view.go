package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var viewCommand = &cobra.Command{
	Use:   "view",
	Short: "View workout plan",
	Long:  "View workout plan",
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
				err = viewCmd(&cfg, user)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		} else {
			err = viewCmd(&cfg, user)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCommand)
}
