package models

import (
	"gorm.io/gorm"
)

type Label struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Color string
}
