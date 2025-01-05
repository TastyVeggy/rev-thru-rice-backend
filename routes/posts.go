package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/comment"
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/post"
	"github.com/labstack/echo/v4"
)

// posts
func PostRoutes(g *echo.Group) {
	g.GET("", post.GetPosts)
	g.GET("/count", post.GetPostCount)
	g.GET("/:id", post.GetPost)
}

// protected/posts
func ProtectedPostRoutes(g *echo.Group) {
	g.PUT("/:id", post.EditPost)
	g.DELETE("/:id", post.DeletePost)
	g.POST("/:id/comments", comment.CreateComment)

	// user cannot create shop individually
	// g.POST("/:id/shops", shop.CreateShop)
}
