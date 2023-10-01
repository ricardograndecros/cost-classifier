package db

import (
	"fmt"

	"cost-classifier/backend/models" // Import your models
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("./transactions.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Automigrate your models
	err = DB.AutoMigrate(&models.Transaction{}, &models.Label{}, &models.User{}, &models.Requisition{})
	if err != nil {
		return err
	}

	fmt.Println("Database connection established")
	return nil
}
