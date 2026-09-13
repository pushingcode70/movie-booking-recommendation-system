package database

import (
	"fmt"
	"log"

	"movie-booking/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// db is the global database connection used throughout the application.
var DB *gorm.DB

// connectdb establishes a connection to the PostgreSQL database.
func ConnectDB() {

	// Build the DSN (Data Source Name) using values loaded from .env.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	// open a connection to PostgreSQL using GORM.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// store the database connection globally.
	DB = db

	log.Println("Database Connected Successfully")
}
