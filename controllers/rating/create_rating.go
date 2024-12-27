package rating

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func CreateRating(c echo.Context) error {
	shopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert shop id parameter to integer")
	}

	rating := new(services.RatingReqDTO)
	if err := c.Bind(rating); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad post request: %v", err))
	}

	userID := c.Get("user").(int)

	ratingRes, err := services.AddRating(rating, shopID, userID)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to insert post: %v", err))
	}

	res := map[string]any{
		"message": "Rating successfully added",
		"rating":  ratingRes,
	}
	return c.JSON(http.StatusOK, res)

}
