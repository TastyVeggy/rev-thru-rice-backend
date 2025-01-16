package post

import (
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetShopReview(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}
	shopReview, err := services.FetchShopReviewByPostID(postID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.String(http.StatusNotFound, "Review not found")

		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, shopReview)
}
