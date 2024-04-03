package main

import (
	// "database/sql"
	// "log"
	// "os"

	"github.com/harljos/gymplanr/cmd"
	// "github.com/harljos/gymplanr/internal/database"
	// "github.com/joho/godotenv"
)

// type config struct {
// 	DB *database.Queries
// }

func main() {
	// godotenv.Load(".env")

	// dbURL := os.Getenv("DB_URL")
	// if dbURL == "" {
	// 	log.Fatal("DB_URL is not found in the environment")
	// }

	// db, err := sql.Open("postgress", dbURL)
	// if err != nil {
	// 	log.Fatal("Can't connect to database:", err)
	// }

	// apiCfg := config{
	// 	DB: database.New(db),
	// }

	cmd.Execute()
}
