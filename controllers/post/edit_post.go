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

	postRes, err := services.UpdatePost(newPost, postID, userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.String(http.StatusUnauthorized, "User cannot change other user's post or post not found")
		} else if err.Error() == "country not part of list" {
			return c.String(http.StatusBadRequest, "Added a country not part of list")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update post: %v", err))
	}

	res := map[string]any{
		"message": "Post updated successfully",
		"post":    postRes,
	}

	return c.JSON(http.StatusOK, res)
}
