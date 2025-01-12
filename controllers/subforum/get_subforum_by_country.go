package subforum

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetSubforumsWithPostCount(c echo.Context) error {
	var countryID *int

	countryIDString := c.QueryParam("country_id")
	if countryIDString == ""{
		countryID = nil
	} else {
		countryIDint, err := strconv.Atoi(countryIDString)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Can't convert country id parameter to integer")
		}
		countryID = &countryIDint
	}

	subforums, err := services.FetchAllSubforumsWithPostCount(countryID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch subforums with post count: %v", err))
	}
	return c.JSON(http.StatusOK, subforums)
}