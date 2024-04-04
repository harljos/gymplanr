package cmd

import (
	"database/sql"
	"errors"
	"os"

	"github.com/harljos/gymplanr/internal/database"
	"github.com/joho/godotenv"
)

type config struct {
	DB *database.Queries
}

func connectToDB() (config, error) {
	godotenv.Load(".env")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return config{}, errors.New("DB_URL is not found in the environment")
	}

	db, err := sql.Open("postgress", dbURL)
	if err != nil {
		return config{}, err
	}

	cfg := config{
		DB: database.New(db),
	}

	return cfg, nil
}
