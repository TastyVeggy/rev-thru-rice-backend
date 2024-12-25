package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type SignupDTO struct {
	Username        string `json:"username" validate:"required,min=3,max=30"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

func Signup(c echo.Context) error {
	var err error
	user := new(SignupDTO)
	if err = c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Signup bad request: %v", err))
	}

	// Backend validation for increased security
	err = utils.Validator.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Validation failed: %s %s %s", err.StructField(), err.Tag(), err.Param()))
		}
	} else {
		log.Println("Validation successful")
	}

	if user.Password != user.ConfirmPassword {
		return c.String(http.StatusBadRequest, "Confirm Password does not match")
	}

	userExists, err := models.Exists("username", user.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to verify if username has been taken: %v", err))
	}
	emailExists, err := models.Exists("email", user.Email)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to verify if email has been taken: %v", err))
	}

	if userExists {
		return c.String(http.StatusBadRequest, "Username has been taken")
	}
	if emailExists {
		return c.String(http.StatusBadRequest, "Email has been taken")
	}

	err = models.AddUser(user.Username, user.Email, user.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error adding new user: %v", err))
	}

	userID, err := models.FindValueByColumn("username", user.Username, "id")
	if err != nil {
		return c.String(http.StatusOK, user.Username+" has been successfully added but unable to generate JWT token")
	}

	err = utils.GenerateJWTandSetCookie(userID, c)
	if err != nil {
		return c.String(http.StatusOK, user.Username+" has been successfully added but unable to generate JWT token")
	}

	return c.String(http.StatusOK, user.Username+" has been successfully added.")
}
