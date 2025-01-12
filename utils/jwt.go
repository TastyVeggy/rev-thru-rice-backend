package utils

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWTandSetCookie(userID int, c echo.Context) error {
	token, err := generateJWT(userID)
	if err != nil {
		return err
	}
	c.SetCookie(&http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path: "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
	return nil
}

func RemoveJWTCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "jwt_token",
		Path: "/",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
}

func generateJWT(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: strconv.Itoa(userID),
			Issuer:    "rev-thru-rice",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func parseJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and return the claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func CheckTokenValidity(c echo.Context) (*Claims,error) {
		var tokenString string

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			if len(authHeader) > 7 && authHeader[:7] == "Bearer "{
				tokenString = authHeader[7:]
			} else {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header")
			}
		}

		if tokenString == ""{
			cookie, err := c.Cookie("jwt_token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
		}


		claims, err := parseJWT(tokenString)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		return claims, nil
}
