package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/labstack/echo/v4"
)

func GetPosts(c echo.Context) error {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	subforum := c.QueryParam("subforum")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	// default to -1, which means just all posts
	if subforum == "" {
		subforum = "-1"
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

	subforumNum, err := strconv.Atoi(subforum)
	if err != nil {
		subforumNum = -1
	}
	offset := (pageNum - 1) * limitNum

	posts, err := models.FetchPosts(limitNum, offset, subforumNum)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch posts: %v", err))
	}

	res := map[string]any{
		"page":  pageNum,
		"limit": limitNum,
		"posts": posts,
	}

	return c.JSON(http.StatusOK, res)

}
