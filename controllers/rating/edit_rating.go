package rating

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func EditRating(c echo.Context) error {
	shopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert shop id parameter to integer")
	}

	newRating := new(services.RatingReqDTO)
	if err := c.Bind(newRating); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad put request: %v", err))
	}

	userID := c.Get("user").(int)

	ratingRes, err := services.UpdateRating(newRating, userID, shopID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.String(http.StatusUnauthorized, "You cannot change other people's rating or rating not found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update rating: %v", err))
	}

	res := map[string]any{
		"message": "Rating updated successfully",
		"rating":  ratingRes,
	}

	return c.JSON(http.StatusOK, res)
}
