package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	ProfilePic string    `json:"profile_pic"`
}

type UserReqDTO struct {
	Username        string `json:"username" validate:"required,min=3,max=30"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
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

func GetUserID(username string) (int, error) {
	var result int

	query := "SELECT id FROM users WHERE username = $1"
	err := db.Pool.QueryRow(context.Background(), query, username).Scan(&result)

	if err != nil {
		if err == pgx.ErrNoRows {
			return result, errors.New("username does not exist")
		}
		return result, err
	}
	return result, nil
}

func GetPassword(username string) (string, error) {
	var result string
	query := "SELECT password FROM users WHERE username = $1"
	err := db.Pool.QueryRow(context.Background(), query, username).Scan(&result)

	if err != nil {
		if err == pgx.ErrNoRows {
			return result, errors.New("username does not exist")
		}
		return result, err
	}
	return result, nil
}

func AddUser(user *UserReqDTO) error {
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
