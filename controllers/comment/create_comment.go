package comment

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func CreateComment(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert post id parameter to integer")
	}

	userID := c.Get("user").(int)

	comment := new(services.CommentReqDTO)
	if err := c.Bind(comment); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad post request: %v", err))
	}

	commentRes, err := services.AddComment(comment, userID, postID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to insert comment: %v", err))
	}

	res := map[string]any{
		"message": "Comment successfully added",
		"comment": commentRes,
	}
	return c.JSON(http.StatusOK, res)
}
