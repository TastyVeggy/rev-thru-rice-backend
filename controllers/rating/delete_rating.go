package rating

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func DeleteRating(c echo.Context) error {
	userID := c.Get("user").(int)
	shopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert shop id parameter to integer")
	}

	RowsDeletedCount, err := services.RemoveRating(shopID, userID)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete rating: %v", err))
	}

	// If no rows affected, it means that current user requesting for deletion does not tally with the user_id associated with the comment
	// Or maybe the comment just doesn't exists
	if RowsDeletedCount == 0 {
		return c.String(http.StatusUnauthorized, "You cannot delete other people's rating or rating not found")
	}

	return c.JSON(http.StatusOK, "Rating deleted successfully")
}
