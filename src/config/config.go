package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB DB is a global variable that holds the connection to the database
func ConnectDB() *gorm.DB {
	errorENV := godotenv.Load()
	if errorENV != nil {
		panic("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, errorDB := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if errorDB != nil {
		panic("Failed to connect to database")
	}
	return db
}

// DisconnectDB is stopping connection to mysql database
func DisconnectDB(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to kill connection from database")
	}
	err = dbSQL.Close()
	if err != nil {
		return
	}
}
