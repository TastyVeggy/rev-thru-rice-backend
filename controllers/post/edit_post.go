package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func EditPost(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}

	newPost := new(services.PostReqDTO)
	if err := c.Bind(newPost); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad put request: %v", err))
	}

	userID := c.Get("user").(int)

	RowsAffectedCount, err := services.EditPost(newPost, postID, userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update post: %v", err))
	}

	// If no rows affected, it means that current user requesting for changing does not tally with the user_id associated with the post
	if RowsAffectedCount == 0 {
		return c.String(http.StatusUnauthorized, "You cannot change other people's post or post not found")
	}
	return c.String(http.StatusOK, "Post updated successfully")

}
