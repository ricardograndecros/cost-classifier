package handlers

import (
	"log"
	"net/http"
	"strconv"

	"cost-classifier/backend/db"
	"cost-classifier/backend/fetcher"
	"cost-classifier/backend/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handle creating a new transaction
func CreateTransaction(c echo.Context) error {
	var transaction models.Transaction
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	result := db.DB.Create(&transaction)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create transaction",
		})
	}

	return c.JSON(http.StatusOK, transaction)
}

// Handle deleting a transaction
func DeleteTransaction(c echo.Context) error {
	id := c.QueryParam("id")

	var transaction models.Transaction
	result := db.DB.Where("id = ?", id).Delete(&transaction)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete transaction",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Transaction deleted successfully",
	})
}

// Handle getting a list of transactions
func GetTransactions(c echo.Context) error {
	var transactions []models.Transaction
	var result *gorm.DB

	countParam := c.QueryParam("count")
	if countParam != "" {
		// If count parameter is provided, retrieve the specified number of transactions
		count, err := strconv.Atoi(countParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid count parameter",
			})
		}
		result = db.DB.Limit(count).Find(&transactions)
	} else {
		// If count parameter is not provided, retrieve all transactions
		result = db.DB.Find(&transactions)
	}
	log.Printf("Count Param: %s", countParam)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve transactions",
		})
	}
	log.Printf("DB Error: %v", result.Error)
	log.Printf("Transactions: %v", transactions)

	return c.JSON(http.StatusOK, transactions)
}

func FetchTransactionsFromBank(c echo.Context) error {
	username := c.Get("userId").(string)
	fetcher.RetrieveAccountTransactions(fetcher.NewNordigenClient(), username)

	return c.JSON(http.StatusCreated, "New transactions processed and stored")
}
