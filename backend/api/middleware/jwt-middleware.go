package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"cost-classifier/backend/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Missing Authorization Header",
			})
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Invalid/Malformed Token",
			})
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.SecretKey), nil
		})

		if err != nil {
			log.Printf("JWT Error: %s", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid/Malformed Token",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Storing user ID in the context
			c.Set("userId", claims["userId"].(string))
			// Pass control to the next handler in the chain
			return next(c)
		} else {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid Token",
			})
		}
	}
}
