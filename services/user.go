package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

type UserReqDTO struct {
	Username        string `json:"username" validate:"required,min=3,max=30"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	//TODO: Profile pic
	// ProfilePic string `json:"string" validate:"required"`
}

type UserResDTO struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	ProfilePic string    `json:"profile_pic"`
}

type LoginReqDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdatePasswordReqDTO struct {
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UpdateUserInfoReqDTO struct {
	Username        string `json:"username" validate:"required,min=3,max=30"`
	Email           string `json:"email" validate:"required,email"`
}

func AddUser(user *UserReqDTO) (UserResDTO, error) {
	var userRes UserResDTO

	err := utils.Validator.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return userRes, fmt.Errorf("validation failed: %v", err)
		}
	}

	if user.Password != user.ConfirmPassword {
		return userRes, errors.New("confirm password does not match")
	}

	// can finally add user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return userRes, fmt.Errorf("error adding new user: %v", err)
	}

	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1,$2,$3)
		RETURNING *
	`

	var password_store string

	err = db.Pool.QueryRow(
		context.Background(),
		query,
		user.Username,
		user.Email,
		string(hashedPassword),
	).Scan(
		&userRes.ID,
		&userRes.Username,
		&userRes.Email,
		&password_store,
		&userRes.CreatedAt,
		&userRes.ProfilePic,
	)
	if err != nil {
		if err.Error() == `duplicate key value violates unique constraint "users_username_key` {
			err = errors.New("username has been taken")
		} else if err.Error() == `duplicate key value violates unique constraint "users_email_key` {
			err = errors.New("email has been taken")
		} else {
			err = fmt.Errorf("error adding new user: %v", err)
		}
	}

	return userRes, err
}

func LoginUser(user *LoginReqDTO) (UserResDTO, error) {
	var userRes UserResDTO

	// First get all data of user by name
	query := "SELECT * FROM users WHERE username = $1"

	var storedHashedPassword string
	err := db.Pool.QueryRow(
		context.Background(),
		query,
		user.Username,
	).Scan(
		&userRes.ID,
		&userRes.Username,
		&userRes.Email,
		&storedHashedPassword,
		&userRes.CreatedAt,
		&userRes.ProfilePic,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return userRes, errors.New("username does not exist")
		} else {
			return userRes, err
		}
	}

	isCorrectPassword := utils.ComparePasswords(storedHashedPassword, user.Password)

	if !isCorrectPassword {
		return userRes, errors.New("password is incorrect")
	}

	return userRes, nil
}

func UpdatePassword(password *UpdatePasswordReqDTO, userID int)(error){
	err := utils.Validator.Struct(password)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("validation failed: %v", err)
		}
	}

	if password.Password != password.ConfirmPassword {
		return  errors.New("confirm password does not match")
	}

	// can finally add user
	hashedPassword, err := utils.HashPassword(password.Password)
	if err != nil {
		return fmt.Errorf("error adding new user: %v", err)
	}

	query := `
		UPDATE users
		SET password=$1
		WHERE user_id=$2
	`

	commandTag, err := db.Pool.Exec(context.Background(), query, hashedPassword, userID)
	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}
	return err
}


func UpdateUserInfo(user *UpdateUserInfoReqDTO, userID int) (UserResDTO, error){
	var userRes UserResDTO

	query := `
			UPDATE users
			SET username=$1, email=$2
			WHERE id=$3
			RETURNING id, username, email, created_at, profile_pic
		`
	err := db.Pool.QueryRow(
			context.Background(),
			query,
			user.Username,
			user.Email,
			userID,
		).Scan(
			&userRes.ID,
			&userRes.Username,
			&userRes.Email,
			&userRes.CreatedAt,
			&userRes.ProfilePic,
		)

	if err != nil {
		return userRes, err
	}

	return userRes, nil
}

func RemoveUser(id int) error {
	query := `
		UPDATE users
		SET username=NULL,email=NULL,password=NULL,created_at=NULL,profile_pic=''
		WHERE id=$1 
	`
	commandTag, err := db.Pool.Exec(context.Background(), query, id)
	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}
	return err
}

func FetchUserByID(id int) (UserResDTO, error) {
	var user UserResDTO

	query := "SELECT username, email FROM users WHERE id = $1"

	err := db.Pool.QueryRow(context.Background(), query, id).Scan(
		&user.Username,
		&user.Email,
	)

	return user, err
}
