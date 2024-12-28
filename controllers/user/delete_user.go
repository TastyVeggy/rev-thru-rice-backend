package user

import (
	"fmt"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

// soft delete (change all information to null)
func DeleteUser(c echo.Context) error {
	userID := c.Get("user").(int)

	RowsDeletedCount, err := services.RemoveUser(userID)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete post: %v", err))
	}

	// 0 rows affected means, user does not exist
	if RowsDeletedCount == 0 {
		return c.String(http.StatusNotFound, "User does not exist")
	}

	return c.JSON(http.StatusOK, "User deleted successfully")
}

