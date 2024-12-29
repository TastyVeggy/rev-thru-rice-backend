package subforum

import (
	"fmt"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/labstack/echo/v4"
)

func GetSubforums(c echo.Context) error {
	subforums, err := models.FetchAllSubforums()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch subforums: %v", err))
	}
	return c.JSON(http.StatusOK, subforums)
}