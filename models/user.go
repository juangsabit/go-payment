package models

import (
	"errors"
	"go-payment/utils/token"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	FullName      string `gorm:"type:varchar(255)" json:"fullname"`
	Username      string `gorm:"type:varchar(255)" binding:"required" json:"username"`
	Password      string `gorm:"type:varchar(255)" binding:"required" json:"password"`
	Email         string `gorm:"type:varchar(255)" json:"email"`
	RoleID        uint   `json:"role_id"`
	Balance       uint   `json:"balance"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	Activitys     []Activity    `gorm:"foreignKey:UserID;references:ID"`
	Transactions  []Transaction `gorm:"foreignKey:SenderUserID;references:ID;"`
	Transactions2 []Transaction `gorm:"foreignKey:ReceiverUserID;references:ID;"`
}

type Activity struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	UserID      uint   `json:"user_id"`
	Information string `gorm:"type:varchar(255)" json:"information"`
	TableName   string `gorm:"type:varchar(50)" json:"table_name"`
	RowID       uint   `json:"row_id"`
	CreatedAt   *time.Time
}

type Role struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"type:varchar(255)" json:"name"`
	Users []User `gorm:"foreignKey:RoleID; references:ID"`
}

func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func LoginCheck(username string, password string) (string, error) {

	var err error
	u := User{}
	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	match := VerifyPassword(password, u.Password)
	if !match {
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password)); err != nil {
			return "", err
		}
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

type UserLogin struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" binding:"required" json:"username"`
	Password string `gorm:"size:255;not null;" binding:"required" json:"password"`
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New(strconv.Itoa(int(uid)) + "User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func SaveActivity(user_id uint, info string, tableName string, rowID uint) {
	obj := Activity{UserID: user_id, Information: info, TableName: tableName, RowID: rowID}
	DB.Debug().Create(&obj)
}

func CheckBalanaceUserByID(uid uint) int {

	userDetail, _ := GetUserByID(uid)
	return int(userDetail.Balance)

}
