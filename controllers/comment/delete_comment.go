package comment

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func DeleteComment(c echo.Context) error {
	userID := c.Get("user").(int)
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Can't convert comment id parameter to integer")
	}

	RowsDeletedCount, err := services.RemoveComment(commentID, userID)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete comment: %v", err))
	}

	// If no rows affected, it means that current user requesting for deletion does not tally with the user_id associated with the comment
	// Or maybe the comment just doesn't exists
	if RowsDeletedCount == 0 {
		return c.String(http.StatusUnauthorized, "You cannot delete other people's comment or comment not found")
	}

	return c.JSON(http.StatusOK, "Comment deleted successfully")
}
