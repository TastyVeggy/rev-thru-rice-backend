package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/country"
	"github.com/labstack/echo/v4"
)

// countries are fixed for now
// countries
func CountryRoutes(g *echo.Group) {
	g.GET("", country.GetCountries)
	// g.GET(":id/post_counts", subforum.GetPostCounts) // id < 1 if looking for the count of posts within each subforum without restriction towards country association of posts
	// // it returns only the post_counts for subforums with more than one post

}
