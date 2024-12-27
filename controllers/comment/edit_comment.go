package comment

import (
	"fmt"
	"net/http"
	"strconv"

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

	RowsAffectedCount, err := services.EditComment(newComment, userID, commentID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update comment: %v", err))
	}

	// If no rows affected, it means that current user requesting for changing does not tally with the user_id associated with the comment
	if RowsAffectedCount == 0 {
		return c.String(http.StatusUnauthorized, "You cannot change other people's comment or comment not found")
	}
	return c.String(http.StatusOK, "Comment updated successfully")

}
