package cmd

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/harljos/gymplanr/internal/database"
)

type Day struct {
	dayName string
	muscles []string
}

func (cfg *config) createDays(days []Day, user database.User) ([]database.Day, error) {
	databaseDays := []database.Day{}

	for _, d := range days {
		day, err := cfg.createDay(d.dayName, user)
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

func (cfg *config) getDayByUser(user database.User, dayName string) (database.Day, error) {
	day, err := cfg.DB.GetDayByNameForUser(context.Background(), database.GetDayByNameForUserParams{
		Name: dayName,
		UserID: user.ID,
	})
	if err != nil {
		return database.Day{}, err
	}

	return day, nil
}

func getWorkoutDays(results map[string]string) ([]Day, error) {
	days, ok := results[daysKey]
	if !ok {
		return []Day{}, errors.New("day result not found")
	}

	minutes, err := strconv.Atoi(results[hoursKey])
	if err != nil {
		return []Day{}, err
	}
	minPerExercise := 7

	fullBodyMuscles := []string{"chest", "lats", "hamstrings", "glutes", "shoulders", "quadriceps", "biceps", "calves", "triceps"}
	upperMuscles := []string{"chest", "lats", "shoulders", "biceps", "triceps", "middle_back"}
	pushMuscles := []string{"chest", "shoulders", "triceps"}
	pullMuscles := []string{"lats", "biceps", "middle_back", "lower_back", "traps"}
	legMuscles := []string{"hamstrings", "glutes", "quadriceps", "calves"}

	if days == "3" {
		return []Day{
			{
				dayName: "Monday",
				muscles: checkOutOfBounds(fullBodyMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Wednesday",
				muscles: checkOutOfBounds(fullBodyMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Friday",
				muscles: checkOutOfBounds(fullBodyMuscles, minutes/minPerExercise),
			},
		}, nil
	}
	if days == "4" {
		return []Day{
			{
				dayName: "Monday",
				muscles: checkOutOfBounds(upperMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Tuesday",
				muscles: checkOutOfBounds(legMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Thursday",
				muscles: checkOutOfBounds(upperMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Friday",
				muscles: checkOutOfBounds(legMuscles, minutes/minPerExercise),
			},
		}, nil
	}
	if days == "5" {
		return []Day{
			{
				dayName: "Monday",
				muscles: checkOutOfBounds(pushMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Tuesday",
				muscles: checkOutOfBounds(pullMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Wednesday",
				muscles: checkOutOfBounds(legMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Friday",
				muscles: checkOutOfBounds(pushMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Saturday",
				muscles: checkOutOfBounds(pullMuscles, minutes/minPerExercise),
			},
		}, nil
	}
	if days == "6" {
		return []Day{
			{
				dayName: "Monday",
				muscles: checkOutOfBounds(pushMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Tuesday",
				muscles: checkOutOfBounds(pullMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Wednesday",
				muscles: checkOutOfBounds(legMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Thursday",
				muscles: checkOutOfBounds(pushMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Friday",
				muscles: checkOutOfBounds(pullMuscles, minutes/minPerExercise),
			},
			{
				dayName: "Saturday",
				muscles: checkOutOfBounds(legMuscles, minutes/minPerExercise),
			},
		}, nil
	}
	return []Day{}, errors.New("no days found")
}

func checkOutOfBounds(muscles []string, exercises int) []string {
	if len(muscles) >= exercises {
		return muscles[0:exercises]
	}

	muscles = append(muscles, muscles...)

	return checkOutOfBounds(muscles, exercises)
}
