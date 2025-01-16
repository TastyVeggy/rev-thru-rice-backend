package user

import (
	"fmt"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func EditPassword(c echo.Context) error {
	userID := c.Get("user").(int)

	newPassword := new(services.UpdatePasswordReqDTO)
	if err := c.Bind(newPassword); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad put request: %v", err))
	}

	err := services.UpdatePassword(newPassword, userID)
	if err != nil {
		if err.Error() == "no row affected" {
			return c.String(http.StatusNotFound, "User cannot be found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update password: %v", err))
	}

	return c.String(http.StatusOK, "Password updated successfully")
}
