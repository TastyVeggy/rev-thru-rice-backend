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

	err = services.RemoveRating(shopID, userID)

	if err != nil {
		if err.Error() == "no row affected"{
			return c.String(http.StatusUnauthorized, "You cannot delete other people's rating or rating not found, or you are the one who posted the shop so you must keep a rating")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete rating: %v", err))
	}

	return c.JSON(http.StatusOK, "Rating deleted successfully")
}
