package post

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetPosts(c echo.Context) error {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	subforumIDString := c.QueryParam("subforum_id")
	userIDString := c.QueryParam("user_id") // posts from this user
	countryIDsString := c.QueryParam("country_ids")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	// default to -1, which means just all posts
	if subforumIDString == "" {
		subforumIDString = "-1"
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

	subforumID, err := strconv.Atoi(subforumIDString)
	if err != nil {
		subforumID = -1
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		userID = -1
	}

	countryIDs := []int{}
	for _, IDstring := range strings.Split(countryIDsString, ",") {
		countryID, err := strconv.Atoi(IDstring)
		if err != nil {
			countryIDs = []int{}
			break
		}
		countryIDs = append(countryIDs, countryID)
	}
	offset := (pageNum - 1) * limitNum

	posts, err := services.FetchPosts(limitNum, offset, subforumID, userID, countryIDs)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch posts: %v", err))
	}

	// res := map[string]any{
	// 	"page":  pageNum,
	// 	"limit": limitNum,
	// 	"posts": posts,
	// }

	return c.JSON(http.StatusOK, posts)

}
