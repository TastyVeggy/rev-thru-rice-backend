package user

import (
	"fmt"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

// soft delete (change all information to null)
func DeleteUser(c echo.Context) error {
	userID := c.Get("user").(int)

	err := services.RemoveUser(userID)

	if err != nil {
		if err.Error() == "no row affected"{
			return c.String(http.StatusNotFound, "User cannot be found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete user: %v", err))
	}
	utils.RemoveJWTCookie(c)

	return c.JSON(http.StatusOK, "User deleted successfully")
}

