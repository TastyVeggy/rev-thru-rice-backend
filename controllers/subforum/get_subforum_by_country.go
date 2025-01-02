package subforum

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/TastyVeggy/rev-thru-rice-backend/services"
// 	"github.com/labstack/echo/v4"
// )

// func GetSubforums(c echo.Context) error {
// 	var countryID *int
// 	var err error

// 	countryIDString := c.Param("id")

// 	if countryIDString != ""{
// 		countryIDint, err := strconv.Atoi(countryIDString)
// 		if err != nil {
// 			return c.String(http.StatusInternalServerError, "Can't convert country id parameter to integer")
// 		}
// 		countryID = &countryIDint
// 	} else {
// 		countryID = nil // no country
// 	}

// 	subforums, err := services.FetchAllSubforums(countryID)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch subforums: %v", err))
// 	}
// 	return c.JSON(http.StatusOK, subforums)
// }