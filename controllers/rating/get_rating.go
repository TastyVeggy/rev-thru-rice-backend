package rating

import (
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetRating(c echo.Context) error {
	shopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert shop id parameter to integer")
	}

	userID := c.Get("user").(int)

	rating, err := services.FetchRatingByShopandUser(shopID, userID)
	if err != nil {
		return c.String(http.StatusNotFound, "Rating not found")
	}
	return c.JSON(http.StatusOK, rating)
}
