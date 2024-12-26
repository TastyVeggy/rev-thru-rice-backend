package post

import (
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/labstack/echo/v4"
)

func GetPost(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}
	post, err := models.FetchPostByID(postID)
	if err != nil {
		return c.String(http.StatusNotFound, "Post not found")
	}
	return c.JSON(http.StatusOK, post)
}
