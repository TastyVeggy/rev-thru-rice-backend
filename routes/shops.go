package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/rating"
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/shop"
	"github.com/labstack/echo/v4"
)

// posts
func ShopRoutes(g *echo.Group) {
	g.GET("/:id", shop.GetShop)
}

// protected/posts
func ProtectedShopRoutes(g *echo.Group) {
	g.PUT("/:id", shop.EditShop)
	// Deletion of shop not exposed to user, if user wants to delete a shop, he will have to delete the post as they are bundled together
	// g.DELETE("/:id", shop.DeleteShop)
	g.GET("/:id/ratings", rating.GetRating)
	g.POST("/:id/ratings", rating.CreateRating)
	g.PUT("/:id/ratings", rating.EditRating)
	g.DELETE("/:id/ratings", rating.DeleteRating)
}
