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

	postRes, err := services.AddPost(post, userID)

	if err != nil {
		if err.Error() == "country not part of list" {
			return c.String(http.StatusBadRequest, "Added a country not part of list")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to insert post: %v", err))
	}

	res := map[string]any{
		"message": "Post successfully added",
		"post":    postRes,
	}
	return c.JSON(http.StatusOK, res)

}
