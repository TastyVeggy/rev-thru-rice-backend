package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/comment"
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/post"
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/shop"
	"github.com/labstack/echo/v4"
)

// posts
func PostRoutes(g *echo.Group) {
	g.GET("", post.GetPosts)
	g.GET("/:id", post.GetPost)
	// g.GET("/:id/comments", comment.GetPostComments)

}

// protected/posts
func ProtectedPostRoutes(g *echo.Group) {
	g.POST("", post.CreatePost)
	g.PUT("/:id", post.EditPost)
	g.DELETE("/:id", post.DeletePost)
	g.POST("/:id/comments", comment.CreateComment)

	g.POST("/:id/shops", shop.CreateShop)
	// g.POST("/shops", post.CreateShopPost) // Creating a shop comes bundled with a post and a rating
}
