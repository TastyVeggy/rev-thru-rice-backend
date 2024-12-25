package auth

import (
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

func Logout(c echo.Context) error {
	utils.RemoveJWTCookie(c)
	return c.String(http.StatusOK, "Logged out successfully")
}
