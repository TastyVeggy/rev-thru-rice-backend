package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// JSON Request body
// - username
// - email
// - password
// - confirm_password

func Signup(c echo.Context) error {
	var err error
	user := new(services.SignupReqDTO)
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

	userExists, err := services.Exists("username", user.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to verify if username has been taken: %v", err))
	}
	emailExists, err := services.Exists("email", user.Email)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to verify if email has been taken: %v", err))
	}

	if userExists {
		return c.String(http.StatusBadRequest, "Username has been taken")
	}
	if emailExists {
		return c.String(http.StatusBadRequest, "Email has been taken")
	}

	err = services.AddUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error adding new user: %v", err))
	}

	userID, err := services.FetchUserIDbyUsername(user.Username)
	if err != nil {
		return c.String(http.StatusOK, user.Username+" has been successfully added but unable to generate JWT token")
	}

	err = utils.GenerateJWTandSetCookie(userID, c)
	if err != nil {
		return c.String(http.StatusOK, user.Username+" has been successfully added but unable to generate JWT token")
	}

	return c.String(http.StatusOK, user.Username+" has been successfully added.")
}
