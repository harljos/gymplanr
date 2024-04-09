package cmd

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/harljos/gymplanr/internal/database"
)

func (cfg *config) createDays(days []string, user database.User) ([]database.Day, error) {
	databaseDays := []database.Day{}

	for _, dayName := range days {
		day, err := cfg.createDay(dayName, user)
		if err != nil {
			return []database.Day{}, err
		}
		databaseDays = append(databaseDays, day)
	}

	return databaseDays, nil
}

func (cfg *config) createDay(dayName string, user database.User) (database.Day, error) {
	day, err := cfg.DB.CreateDay(context.Background(), database.CreateDayParams{
		ID:        uuid.New(),
		Name:      dayName,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return database.Day{}, err
	}

	return day, nil
}
