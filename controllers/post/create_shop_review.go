package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func CreateShopReview(c echo.Context) error {
	subforumID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert subforum id parameter to integer")
	}
	shopReview := new(services.ShopReviewReqDTO)
	if err := c.Bind(shopReview); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad post request: %v", err))
	}

	userID := c.Get("user").(int)

	shopReviewRes, err := services.AddShopReview(shopReview, userID, subforumID)

	if err != nil {
		if err.Error() == "cannot add shop review to non shop review subforums" || err.Error() == "should not have any countries in shop post request, country is determined via lat long"{
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to create shop post: %v", err))
	}

	res := map[string]any{
		"message":   "Shop post successfully added",
		"review": shopReviewRes,
	}
	return c.JSON(http.StatusOK, res)

}
