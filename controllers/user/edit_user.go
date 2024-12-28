package user

// TODO: implement edit user info and password
// func EditUser(c echo.Context) error {
// 	userID := c.Get("user").(int)

// 	newUser := new(services.UserReqDTO)
// 	if err := c.Bind(newUser); err != nil {
// 		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad put request: %v", err))
// 	}

// 	userRes, err := services.UpdateUser(newUser, userID)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "no rows in result set") {
// 			return c.String(http.StatusNotFound, "User cannot be found")
// 		}
// 		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to update user info: %v", err))
// 	}

// 	res := map[string]any{
// 		"message": "User info updated successfully",
// 		"user":    userRes,
// 	}

// 	return c.JSON(http.StatusOK, res)
// }