package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	ProductName string `gorm:"type:varchar(255)" binding:"required" json:"product_name"`
	Description string `gorm:"type:varchar(255)" binding:"required" json:"description"`
	Price       int64  `binding:"required" json:"price"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
