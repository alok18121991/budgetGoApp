package configs

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var JwtSecret = []byte("your_secret_key")
var JwtRefreshSecret = []byte("your_secret_key_refresh")

type Claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

type TokenClaims struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Exclude the login route from JWT authentication
		if c.Path() == "/auth/user" || c.Path() == "/auth/refresh" {
			return next(c)
		}

		// Retrieve the token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
		}

		// Token format should be "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token format")
		}
		tokenString := tokenParts[1]

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtSecret, nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		// Verify token validity
		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		// Token is valid, proceed to the next handler
		return next(c)
	}
}
