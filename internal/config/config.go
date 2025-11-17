package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds all application configurations.
type AppConfig struct {
	DatabaseURL     string
	JWTSecret       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

// LoadConfig loads configuration from environment variables or .env file.
func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, assuming environment variables are set.")
	}

	accessTokenExpHoursStr := os.Getenv("ACCESS_TOKEN_EXP_HOURS")
	if accessTokenExpHoursStr == "" {
		accessTokenExpHoursStr = "1" // Default to 1 hour
	}
	accessTokenExpHours, err := strconv.Atoi(accessTokenExpHoursStr)
	if err != nil {
		log.Fatalf("Invalid ACCESS_TOKEN_EXP_HOURS: %v", err)
	}

	refreshTokenExpDaysStr := os.Getenv("REFRESH_TOKEN_EXP_DAYS")
	if refreshTokenExpDaysStr == "" {
		refreshTokenExpDaysStr = "7" // Default to 7 days
	}
	refreshTokenExpDays, err := strconv.Atoi(refreshTokenExpDaysStr)
	if err != nil {
		log.Fatalf("Invalid REFRESH_TOKEN_EXP_DAYS: %v", err)
	}

	return &AppConfig{
		DatabaseURL:     getEnv("DATABASE_URL", ""),                // We'll construct this from DB_* variables
		JWTSecret:       getEnv("JWT_SECRET", "supersecretjwtkey"), // Replace with a strong, random key in production
		AccessTokenExp:  time.Duration(accessTokenExpHours) * time.Hour,
		RefreshTokenExp: time.Duration(refreshTokenExpDays) * 24 * time.Hour,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
