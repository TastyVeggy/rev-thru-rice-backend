package post

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	postRes, err := services.EditPost(newPost, postID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.String(http.StatusUnauthorized, "You cannot change other people's post or post not found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update post: %v", err))
	}

	res := map[string]any{
		"message": "Post updated successfully",
		"post":    postRes,
	}

	return c.JSON(http.StatusOK, res)
}
