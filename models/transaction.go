package models

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	SenderUserID   uint   `json:"sender_user_id"`
	ReceiverUserID uint   `json:"receiver_user_id"`
	ChannelID      uint   `json:"channel_id"`
	Amount         uint   `json:"amount"`
	Status         string `gorm:"type:varchar(255)" json:"status"`
	CreatedAt      *time.Time
}

type Channel struct {
	gorm.Model
	ID           uint          `gorm:"primaryKey;autoIncrement"`
	Name         string        `gorm:"type:varchar(255)" json:"name"`
	Transactions []Transaction `gorm:"foreignKey:ChannelID;references:ID"`
}

func GetTransactionById(id uint) (Transaction, error) {

	var trx Transaction

	if err := DB.First(&trx, id).Error; err != nil {
		return trx, errors.New(strconv.Itoa(int(id)) + "Transaction not found!")
	}

	return trx, nil

}
