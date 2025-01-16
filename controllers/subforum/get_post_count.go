package subforum

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetPostCounts(c echo.Context) error {
	countryID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert country id parameter to integer")
	}

	subforums, err := services.FetchSubforumPostCountsbyCountryID(countryID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch post counts: %v", err))
	}
	return c.JSON(http.StatusOK, subforums)
}
