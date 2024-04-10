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

func (cfg *config) getDaysByUser(user database.User) ([]database.Day, error) {
	days, err := cfg.DB.GetDaysByUser(context.Background(), user.ID)
	if err != nil {
		return []database.Day{}, err
	}

	return days, nil
}

func (cfg *config) deleteDays(user database.User) error {
	err := cfg.DB.DeleteDays(context.Background(), user.ID)
	if err != nil {
		return err
	}

	return nil
}
