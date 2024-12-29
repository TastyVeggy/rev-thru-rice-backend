package shop

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

// just an internal method for now. Due to the design of the forum making shop tied to a post
func DeleteShop(c echo.Context) error {
	userID := c.Get("user").(int)
	shopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert shop id parameter to integer")
	}

	err = services.RemoveShop(shopID, userID)

	if err != nil {
		if err.Error() == "no row affected" {
			return c.String(http.StatusUnauthorized, "You cannot delete other people's shop or shop not found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete shop: %v", err))
	}

	return c.JSON(http.StatusOK, "Rating deleted successfully")
}
