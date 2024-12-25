package models

import (
	"context"
	"fmt"
	"time"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID         int    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	ProfilePic string    `json:"profile_pic"`
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

func FindValueByColumn(searchCol string, searchVal string, returnCol string) (string, error) {

	var result string

	query := fmt.Sprintf("SELECT %s FROM users where %s = $1", returnCol, searchCol)
	err := db.Pool.QueryRow(context.Background(), query, searchVal).Scan(&result)

	if err != nil {
		if err == pgx.ErrNoRows {
			return result, fmt.Errorf("no record for searchVal: %v", err)

		}
		return result, err
	}
	return result, nil
}

func AddUser(username string, email string, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(context.Background(), "INSERT INTO users(username, email, password) VALUES ($1,$2,$3)", username, email, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}
