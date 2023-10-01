package main

import (
	"cost-classifier/backend/api" // Import your api package
	"cost-classifier/backend/config"
	"cost-classifier/backend/db" // Import your database package
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

func main() {
	// initialize config global variables
	config.InitConfig()

	time.Sleep(time.Second * 3)

	fmt.Println("SECRET KEY: ", config.AppConfig.SecretKey)
	// Initialize the database
	if err := db.InitDatabase(); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:1234", "https://localhost:1234"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	//e.Use(middleware.CORS())

	// Define the routes
	api.DefineRoutes(e)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
	e.Logger.Fatal(e.StartTLS(":443", "cert/server.crt", "cert/server.key"))
}
