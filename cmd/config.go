package cmd

import (
	"database/sql"
	"errors"
	"os"

	"github.com/harljos/gymplanr/internal/database"
	"github.com/harljos/gymplanr/internal/exerciseapi"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type config struct {
	DB             *database.Queries
	exerciseClient exerciseapi.Client
}

func connectToDB() (config, error) {
	godotenv.Load(".env")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return config{}, errors.New("DB_URL is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return config{}, err
	}

	cfg := config{
		DB: database.New(db),
		exerciseClient: exerciseapi.NewClient(),
	}

	return cfg, nil
}
