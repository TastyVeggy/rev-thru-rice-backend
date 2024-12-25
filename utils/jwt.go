package utils

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GenerateJWTandSetCookie(userID string, c echo.Context) error {
	token, err := generateJWT(userID)
	if err != nil {
		return err
	}
	c.SetCookie(&http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   os.Getenv("GO_ENV") != "development",
	})
	return nil
}

func RemoveJWTCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   os.Getenv("GO_ENV") != "development",
	})
}

func generateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "my_server",
			"sub": userID,
		})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
