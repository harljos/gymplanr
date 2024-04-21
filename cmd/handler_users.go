package cmd

import (
	"context"
	"database/sql"
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

	sqlUsername := sql.NullString{
		String: username,
		Valid: true,
	}

	sqlPassword := sql.NullString{
		String: string(hashed),
		Valid: true,
	}

	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Username:  sqlUsername,
		Password:  sqlPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return database.User{}, err
	}

	return user, nil
}

func (cfg *config) loginUserHandler(username, password string) (database.User, error) {
	sqlUsername := sql.NullString{
		String: username,
		Valid: true,
	}

	user, err := cfg.DB.GetUserByUsername(context.Background(), sqlUsername)
	if err != nil {
		return database.User{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
	if err != nil {
		return database.User{}, errors.New("incorrect password")
	}

	return user, nil
}
