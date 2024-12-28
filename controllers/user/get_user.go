package user

import (
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert user id parameter to integer")
	}
	user, err := services.FetchUserByID(userID)
	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}