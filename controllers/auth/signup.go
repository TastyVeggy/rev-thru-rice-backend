package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

// JSON Request body
// - username
// - email
// - password
// - confirm_password

func Signup(c echo.Context) error {
	user := new(services.UserReqDTO)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Signup bad request: %v", err))
	}

	userRes, err := services.AddUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "error adding new user:") {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}

	err = utils.GenerateJWTandSetCookie(userRes.ID, c)
	var res map[string]any
	if err != nil {
		res = map[string]any{
			"message": fmt.Sprintf("Successfully added user but unable to generate JWT token or set cookie due to %v", err),
			"user":    userRes,
		}
	} else {
		res = map[string]any{
			"message": "Successfully added user, generated JWT token and set cookie",
			"user":    userRes,
		}
	}

	return c.JSON(http.StatusOK, res)
}
