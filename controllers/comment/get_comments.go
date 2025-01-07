package comment

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetComments(c echo.Context) error {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	postIDString := c.QueryParam("post_id")
	userIDString := c.QueryParam("user_id") // comments from this user

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	// default to -1, which means just all comments
	if postIDString == "" {
		postIDString = "-1"
	}
	if userIDString == "" {
		userIDString = "-1"
	}

	// Convert query params to integers
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}

	postID, err := strconv.Atoi(postIDString)
	if err != nil {
		postID = -1
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		userID = -1
	}

	offset := (pageNum - 1) * limitNum

	comments, err := services.FetchComments(limitNum, offset, postID, userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch comments: %v", err))
	}

	return c.JSON(http.StatusOK, comments)

}
