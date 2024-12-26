package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/labstack/echo/v4"
)

// does not truly delete post, just set everything to null except id so comments remain
func DeletePost(c echo.Context) error {
	userID := c.Get("user").(int)
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}

	RowsDeletedCount, err := models.RemovePost(postID, userID)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete post: %v", err))
	}

	// If no rows affected, it means that current user requesting for deletion does not tally with the user_id associated with the post
	// Or maybe the post just doesn't exists
	if RowsDeletedCount == 0 {
		return c.String(http.StatusUnauthorized, "You cannot delete other people's post or post not found")
	}

	return c.JSON(http.StatusOK, "Post deleted successfully")
}