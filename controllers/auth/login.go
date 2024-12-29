package auth

import (
	"fmt"
	"net/http"

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

	userRes, err := services.LoginUser(user)
	if err != nil {
		if err.Error() == "username does not exist" || err.Error() == "password is incorrect" {
			return c.String(http.StatusUnauthorized, "Invalid username or password")
		} else {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to verify credentials, %v", err))
		}
	}

	err = utils.GenerateJWTandSetCookie(userRes.ID, c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Login credentials are correct but unable to generate JWT token or set cookie: %v", err))
	}

	res := map[string]any{
		"message": "Successfully logged in",
		"user":    userRes,
	}

	return c.JSON(http.StatusOK, res)
}
