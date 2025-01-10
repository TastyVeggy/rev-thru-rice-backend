package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/user"
	"github.com/labstack/echo/v4"
)

func UserRoutes(g *echo.Group) {
	g.GET("/:id", user.GetUser)
}

func ProtectedUserRoutes(g *echo.Group) {
	g.PUT("/info", user.EditUserInfo)
	g.PUT("/password", user.EditPassword)
	g.DELETE("", user.DeleteUser)
}
