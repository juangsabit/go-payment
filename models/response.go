package models

import "time"

type ResponseUser struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResponseTransaction struct {
	ID       uint   `json:"id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   uint   `json:"amount"`
	Status   string `json:"status"`
}

type ResponseActivity struct {
	ID          uint       `json:"id"`
	Username    string     `gorm:"type:varchar(255)" json:"username"`
	Information string     `gorm:"type:varchar(255)" json:"information"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type ResponseProduct struct {
	ID          uint               `json:"id"`
	ProductName string             `json:"product_name"`
	Description string             `json:"description"`
	Price       int64              `json:"price"`
	CreatedAt   *time.Time         `json:"createdAt"`
	Activity    []ResponseActivity `json:"activityProduct"`
}
