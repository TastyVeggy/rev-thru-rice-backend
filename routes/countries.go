package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/country"
	"github.com/labstack/echo/v4"
)

// countries are fixed for now
// countries
func CountryRoutes(g *echo.Group) {
	g.GET("", country.GetCountries)
}