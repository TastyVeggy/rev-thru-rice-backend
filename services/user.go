package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/jackc/pgx/v5"
)

type SignupReqDTO struct {
	Username        string `json:"username" validate:"required,min=3,max=30"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginReqDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Exists(field string, value string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM users WHERE %s = $1)", field)
	err := db.Pool.QueryRow(context.Background(), query, value).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, err
}

func FetchUserIDbyUsername(username string) (int, error) {
	var UserID int

	query := "SELECT id FROM users WHERE username = $1"
	err := db.Pool.QueryRow(context.Background(), query, username).Scan(&UserID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return UserID, errors.New("username does not exist")
		}
		return UserID, err
	}
	return UserID, nil
}

func FetchPasswordbyUsername(username string) (string, error) {
	var Password string
	query := "SELECT password FROM users WHERE username = $1"
	err := db.Pool.QueryRow(context.Background(), query, username).Scan(&Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			return Password, errors.New("username does not exist")
		}
		return Password, err
	}
	return Password, nil
}

func AddUser(user *SignupReqDTO) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(context.Background(), "INSERT INTO users(username, email, password) VALUES ($1,$2,$3)", user.Username, user.Email, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}
