package models

import (
	"time"
)

type Transaction struct {
	TransactionID string    `gorm:"primary_key"`
	Title         string    `gorm:"not null"`
	Amount        float64   `gorm:"not null"`
	Currency      string    `gorm:"not null"`
	AccountId     string    `gorm:"not null"`
	AccountIban   string    `gorm:"not null"`
	Date          time.Time `gorm:"not null"`
}
