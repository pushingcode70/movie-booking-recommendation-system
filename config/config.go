package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RazorpayKeyID     string
	RazorpayKeySecret string

	SMTPHost     string
	SMTPPort     int
	SMTPEmail    string
	SMTPPassword string

	TMDBAPIKey string

	JWTSecret string

	EmbeddingServiceURL string

	
}

var AppConfig Config

func LoadConfig() {

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT")) //strconv.Atoi() stands for ASCII to Integer.
	AppConfig = Config{

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		RazorpayKeyID:     os.Getenv("RAZORPAY_KEY_ID"),
		RazorpayKeySecret: os.Getenv("RAZORPAY_KEY_SECRET"),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     port,
		SMTPEmail:    os.Getenv("SMTP_EMAIL"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),

		TMDBAPIKey: os.Getenv("TMDB_API_KEY"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		EmbeddingServiceURL: os.Getenv("EMBEDDING_SERVICE_URL"),

		
	}
}
