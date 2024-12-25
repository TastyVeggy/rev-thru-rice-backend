package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

type (
	LoginDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

func Login(c echo.Context) error {
	var err error
	user := new(LoginDTO)
	if err = c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Login bad request: %v", err))
	}

	userExists, err := models.Exists("username", user.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to verify if user exists: %v", err))
	}
	storedHashedPassword, err := models.FindValueByColumn("username", user.Username, "password")

	// Do this so that when username given does not exist, it will not throw an erroneous internal server error
	if err != nil && !strings.Contains(err.Error(), "no record for searchVal") {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to obtain stored password: %v", err))
	}
	correctPassword := utils.ComparePasswords(storedHashedPassword, user.Password)

	if !userExists || !correctPassword {
		return c.String(http.StatusUnauthorized, "Invalid username or passsword")
	}

	userID, err := models.FindValueByColumn("username", user.Username, "id")
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("Valid login credential but unable to obtain userid to generate JWT:  %v", err))
	}

	err = utils.GenerateJWTandSetCookie(userID, c)
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("Valid login credentials but unable to generate JWT and set cookie: %v", err))
	}

	return c.String(http.StatusOK, "Successfully logged in")
}
