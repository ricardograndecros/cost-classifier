package handlers

import (
	"log"
	"net/http"
	"time"

	"cost-classifier/backend/config"
	"cost-classifier/backend/db"
	"cost-classifier/backend/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	log.Printf("Username: %s, Password: %s", username, password)

	// Query the database for the user
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error

	if err != nil {
		log.Printf("User not found, %v", err)
		return echo.ErrUnauthorized
	}

	// Compare the hashed password in the database with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		// handle error (e.g., password does not match)
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = username
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.AppConfig.SecretKey))
	if err != nil {
		log.Printf("Error signing token: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Something went wrong",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
