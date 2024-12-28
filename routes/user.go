package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/user"
	"github.com/labstack/echo/v4"
)

func UserRoutes(g *echo.Group){
	g.GET("/:id", user.GetUser)
}


func ProtectedUserRoutes(g *echo.Group) {
	// TODO: edit user password and info and profile pic
	// g.PUT("", user.EditUser)
	// soft delete
	g.DELETE("", user.DeleteUser)
}