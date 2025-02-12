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
	g.GET("_with_post_count", subforum.GetSubforumsWithPostCount)
}

// protected/subforums
func ProtectedSubforumRoutes(g *echo.Group) {
	g.POST("/:id/posts", post.CreatePost)
	g.POST("/:id/shop_review", post.CreateShopReview) // Creating a review comes bundled with a post, shop and a rating
}
