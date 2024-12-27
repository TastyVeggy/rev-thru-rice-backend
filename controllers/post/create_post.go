package post

import (
	"fmt"
	"net/http"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func CreatePost(c echo.Context) error {

	post := new(services.PostReqDTO)
	if err := c.Bind(post); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad post request: %v", err))
	}

	userID := c.Get("user").(int)

	err := services.AddPost(post, userID)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to insert post: %v", err))
	}
	return c.String(http.StatusOK, "Post successfully added")

}
