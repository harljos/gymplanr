package main

import (
	"log"
	"os"

	"github.com/harljos/gymplanr/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	cmd.Execute()
}
