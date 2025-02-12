package post

import (
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetShopRating(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}

	userID := c.Get("user").(int)

	rating, err := services.FetchRatingByPostandUser(postID, userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.String(http.StatusNotFound, "Rating not found")

		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, rating)
}
