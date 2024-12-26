package routes

import (
	"github.com/TastyVeggy/rev-thru-rice-backend/controllers/comment"
	"github.com/labstack/echo/v4"
)

func CommentRoutes(g *echo.Group){
	g.GET("", comment.GetComments)
}

func ProtectedCommentRoutes(g *echo.Group) {
	g.PUT("/:id", comment.EditComment)
	g.DELETE("/:id",comment.DeleteComment)
}