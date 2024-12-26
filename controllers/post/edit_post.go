package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/labstack/echo/v4"
)

// JSON request body
// - subforum_id
// - title
// - content
func EditPost(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}

	newPost := new(models.PostReqDTO)
	if err := c.Bind(newPost); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad put request: %v", err))
	}

	newPost.UserID = c.Get("user").(int)

	RowsAffectedCount, err := models.EditPost(postID, newPost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update post: %v", err))
	}

	// If no rows affected, it means that current user requesting for changing does not tally with the user_id associated with the post
	if RowsAffectedCount == 0 {
		return c.String(http.StatusUnauthorized, "You cannot change other people's post")
	}
	return c.String(http.StatusOK, "Post updated successfully")

}
