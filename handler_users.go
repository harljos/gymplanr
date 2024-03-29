package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/harljos/gymplanr/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUserHandler(username, password string) (database.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, err
	}

	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Username:  username,
		Password:  string(hashed),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return database.User{}, err
	}

	return user, nil
}
