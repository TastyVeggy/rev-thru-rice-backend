package middleware

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

func JWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.CheckTokenValidity(c)
		if err != nil {
			return err
		}
		c.Set("user", claims.UserID)
		return next(c)
	}
}