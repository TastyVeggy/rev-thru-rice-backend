package post

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/services"
	"github.com/labstack/echo/v4"
)

func GetPostCount(c echo.Context) error {
	subforumIDString := c.QueryParam("subforum_id")
	userIDString := c.QueryParam("user_id") // posts from this user
	countryIDsString := c.QueryParam("country_ids")

	// default to -1, which means just all posts
	if subforumIDString == "" {
		subforumIDString = "-1"
	}
	if userIDString == "" {
		userIDString = "-1"
	}

	// Convert query params to integers
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
	totalPostCount, err := services.FetchPostCount(subforumID, userID, countryIDs)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to fetch post count: %v", err))
	}

	return c.String(http.StatusOK, strconv.Itoa(totalPostCount))

}