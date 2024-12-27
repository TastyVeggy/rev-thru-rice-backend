package shop

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func CreateShop(c echo.Context) error {
	var err error
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}

	userID := c.Get("user").(int)

	shop := new(services.ShopReqDTO)
	if err := c.Bind(shop); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad post request: %v", err))
	}

	err = services.AddShop(shop, userID, postID)

	if err != nil {
		if strings.Contains(err.Error(), "error in getting location"){
			if strings.Contains(err.Error(), "no results found for given coordinates"){
				return c.String(http.StatusBadRequest, fmt.Sprintf("lat long given invalid: %v", err))
			} else {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to reverse geocode: %v", err))
			}
		} else if strings.Contains(err.Error(), "country not part of list"){
			return c.String(http.StatusBadRequest, fmt.Sprintf("Shop location not within country list: %v", err))
		} else {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to add shop, %v", err))
		}
	}

	return c.String(http.StatusOK, "yeah")
}