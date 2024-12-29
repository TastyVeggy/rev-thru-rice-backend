package shop

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func EditShop(c echo.Context) error {
	shopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert shop_id parameter to integer")
	}

	newShop := new(services.ShopReqDTO)
	if err := c.Bind(newShop); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad put request: %v", err))
	}

	userID := c.Get("user").(int)

	shopRes, err := services.UpdateShop(newShop, userID, shopID)
	if err != nil {
		if err.Error() == "no entry in shop" {
			return c.String(http.StatusUnauthorized, "User cannot change other user's shop or shop not found")
		} else if err.Error() == "country not part of list" {
			return c.String(http.StatusBadRequest, "shop not in country part of list")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update shop: %v", err))
	}

	res := map[string]any{
		"message": "Shop updated successfully",
		"shop":    shopRes,
	}

	return c.JSON(http.StatusOK, res)
}
