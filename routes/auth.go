package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/auth"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(g *echo.Group) {
	g.POST("/signup", auth.Signup)
	g.POST("/login", auth.Login)
	g.POST("/logout", auth.Logout)
	g.POST("/verify_token", auth.VerifyToken)
}
