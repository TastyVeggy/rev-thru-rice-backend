package post

import (
	"fmt"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/labstack/echo/v4"
)

// JSON request body
// - subforum_id
// - title
// - content
// TODO: validation stuff for the post

func CreatePost(c echo.Context) error {

	post := new(models.PostReqDTO)
	if err := c.Bind(post); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad post request: %v", err))
	}

	post.UserID = c.Get("user").(int)

	err := models.AddPost(post)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to insert post: %v", err))
	}
	return c.String(http.StatusOK, "Post successfully added")

}
