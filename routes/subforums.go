package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/post"
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/subforum"
	"github.com/labstack/echo/v4"
)

// subforums are fixed for now
// subforums
func SubforumRoutes(g *echo.Group) {
	g.GET("", subforum.GetSubforums)
}

// protected/subforums
func ProtectedSubforumRoutes(g *echo.Group) {
	g.POST("/:id/posts", post.CreatePost)
	g.POST("/:id/shop_posts", post.CreateShopPost) // Creating a shop comes bundled with a post and a rating
}
