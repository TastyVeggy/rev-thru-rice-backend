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

	err = services.RemoveComment(commentID, userID)

	if err != nil {
		if err.Error() == "no row affected"{
			return c.String(http.StatusUnauthorized, "You cannot delete other people's comment or comment not found")
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete comment: %v", err))
	}

	return c.JSON(http.StatusOK, "Comment deleted successfully")
}
