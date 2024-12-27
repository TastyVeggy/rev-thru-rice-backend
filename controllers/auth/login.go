package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	var err error
	user := new(services.LoginReqDTO)
	if err = c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Login bad request: %v", err))
	}

	userExists, err := services.Exists("username", user.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to verify if user exists: %v", err))
	}

	storedHashedPassword, err := services.FetchPasswordbyUsername(user.Username)
	// Do this so that when username given does not exist, it will not throw an erroneous internal server error
	if err != nil && !strings.Contains(err.Error(), "username does not exist") {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to obtain stored password: %v", err))
	}

	isCorrectPassword := utils.ComparePasswords(storedHashedPassword, user.Password)

	if !userExists || !isCorrectPassword {
		return c.String(http.StatusUnauthorized, "Invalid username or passsword")
	}

	userID, err := services.FetchUserIDbyUsername(user.Username)
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("Valid login credential but unable to obtain userid to generate JWT:  %v", err))
	}

	err = utils.GenerateJWTandSetCookie(userID, c)
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("Valid login credentials but unable to generate JWT and set cookie: %v", err))
	}

	return c.String(http.StatusOK, "Successfully logged in")
}
