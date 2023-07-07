package models

import (
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := "host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Asia/Jakarta"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&Product{}, &Role{}, &User{}, &Activity{}, &Channel{}, &Transaction{})

	DB = database
}

func StrToUint(str string) uint {
	intType, _ := strconv.Atoi(str)
	return uint(intType)
}

func IntToString(param int) string {
	// intType, _ := strconv.Atoi(str)
	// return uint(intType)
	return strconv.FormatUint(uint64(param), 10)
}
