package middleware

import (
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/labstack/echo/v4"
)

func JWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
		}

		tokenString := authHeader[7:]

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		c.Set("user", claims.UserID)

		return next(c)
	}
}
