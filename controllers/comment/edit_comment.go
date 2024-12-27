package comment

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func EditComment(c echo.Context) error {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert comment id parameter to integer")
	}

	newComment := new(services.CommentReqDTO)
	if err := c.Bind(newComment); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad put request: %v", err))
	}

	userID := c.Get("user").(int)

	commentRes, err := services.UpdateComment(newComment, userID, commentID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.String(http.StatusUnauthorized, "You cannot change other people's comment or comment not found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update comment: %v", err))
	}

	res := map[string]any{
		"message": "Comment updated successfully",
		"comment": commentRes,
	}

	return c.JSON(http.StatusOK, res)
}
