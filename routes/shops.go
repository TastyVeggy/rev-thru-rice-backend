package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/rating"
	"github.com/labstack/echo/v4"
)

// posts
// func ShopRoutes(g *echo.Group) {
// 	g.GET("", post.GetPosts)
// 	g.GET("/:id", post.GetPost)
// 	// g.GET("/:id/comments", comment.GetPostComments)

// }

// protected/posts
func ProtectedShopRoutes(g *echo.Group) {
	g.GET("/:id/ratings", rating.GetRating)
	g.POST("/:id/ratings", rating.CreateRating)
	g.PUT("/:id/ratings", rating.EditRating)
	g.DELETE("/:id/ratings", rating.DeleteRating)
}
