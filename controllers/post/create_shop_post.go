package post

import (
	"fmt"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func CreateShopPost(c echo.Context) error {

	shopPost := new(services.ShopPostReqDTO)
	if err := c.Bind(shopPost); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad post request: %v", err))
	}

	userID := c.Get("user").(int)

	shopPostRes, err := services.AddShopPost(shopPost, userID)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to create shop post: %v", err))
	}


	res := map[string]any{
		"message": "Shop post successfully added",
		"shop_post":    shopPostRes,
	}
	return c.JSON(http.StatusOK, res)

}