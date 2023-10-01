package api

import (
	"cost-classifier/backend/api/handlers"
	"cost-classifier/backend/api/middleware"
	"github.com/labstack/echo/v4"
)

func DefineRoutes(e *echo.Echo) {
	// Secured Routes
	secured := e.Group("")
	secured.Use(middleware.IsAuthenticated)

	e.POST("/login", handlers.Login)

	// Labels routes
	secured.POST("/labels/new", handlers.CreateNewLabel)
	secured.PUT("/labels/edit", handlers.EditLabel)
	secured.DELETE("/labels/delete", handlers.DeleteLabel)
	secured.GET("/labels", handlers.GetLabels)

	// Transactions routes
	secured.POST("/transactions/new", handlers.CreateTransaction)
	secured.DELETE("/transactions/delete", handlers.DeleteTransaction)
	secured.GET("/transactions", handlers.GetTransactions)
	secured.POST("/transactions/fetch", handlers.FetchTransactionsFromBank)
}
