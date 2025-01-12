package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func EditUserInfo(c echo.Context) error {
	userID := c.Get("user").(int)

	newUserInfoReq := new(services.UpdateUserInfoReqDTO)
	if err := c.Bind(newUserInfoReq); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad put request: %v", err))
	}

	userRes, err := services.UpdateUserInfo(newUserInfoReq, userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return c.String(http.StatusNotFound, "User cannot be found")
		} else if strings.Contains(err.Error(), "error updating new user:"){
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusBadRequest, err.Error())
	}

	res := map[string]any{
		"message": "User info updated successfully",
		"user":    userRes,
	}

	return c.JSON(http.StatusOK, res)
}
