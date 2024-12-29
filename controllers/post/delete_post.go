package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func DeletePost(c echo.Context) error {
	userID := c.Get("user").(int)
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}

	err = services.RemovePost(postID, userID)

	if err != nil {
		if err.Error() == "no row affected" {
			return c.String(http.StatusUnauthorized, "You cannot delete other people's post or post not found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete post: %v", err))
	}

	return c.JSON(http.StatusOK, "Post deleted successfully")
}
