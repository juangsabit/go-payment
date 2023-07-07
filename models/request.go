package models

import (
	"gorm.io/gorm"
)

type ApproveTransactionRequest struct {
	gorm.Model
	ID     uint   `json:"id"`
	Status string `gorm:"type:varchar(255)" json:"status"`
}
