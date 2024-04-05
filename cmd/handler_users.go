package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/harljos/gymplanr/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *config) createUserHandler(username, password string) (database.User, error) {
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

func (cfg *config) loginUserHandler(username, password string) (database.User, error) {
	user, err := cfg.DB.GetUserByUsername(context.Background(), username)
	if err != nil {
		return database.User{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return database.User{}, errors.New("incorrect password")
	}

	return user, nil
}
