package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

func VerifyToken(c echo.Context) error {
	claims, err := utils.CheckTokenValidity(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, fmt.Sprintf("invalid or expired token: %v", err))
	}

	user, err := services.FetchUserByID(claims.UserID)
	if (err != nil){
		if strings.Contains(err.Error(), "no rows in result"){
			return c.String(http.StatusNotFound, "user not found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch user: %v", err))
	}

	res := map[string]any{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
	}
	return c.JSON(http.StatusOK, res)
}