package database

import (
	"fmt"
	"log"

	"movie-booking/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// db is the global database connection used throughout the application
var DB *gorm.DB

// connectDB establishes a connection to the postgresql database
func ConnectDB() {

	// build dsn using values loaded from .env
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	// open connection to postgresql using gorm
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// store database connection globally
	DB = db

	log.Println("Database Connected Successfully")
}
