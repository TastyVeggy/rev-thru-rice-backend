package auth

import (
	"fmt"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

func VerifyToken(c echo.Context) error {
	claims, err := utils.CheckTokenValidity(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, fmt.Sprintf("invalid or expired token: %v", err))
	}
	res := map[string]any{
		"id": claims.UserID,
		"username": claims.Username,
		"email": claims.Email,
	}
	return c.JSON(http.StatusOK, res)
}