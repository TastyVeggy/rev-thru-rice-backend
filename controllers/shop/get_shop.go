package shop

import (
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetShop(c echo.Context) error {
	shopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert shop id parameter to integer")
	}
	shop, err := services.FetchShopByID(shopID)
	if err != nil {
		if err.Error() == "no rows in result set"{
			return c.String(http.StatusNotFound, "Shop not found")
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, shop)
}
