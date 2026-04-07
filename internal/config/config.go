package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	JWTAccessSecret  string
	JWTRefreshSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Port:             getEnv("PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPass:           getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "userapi"),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "your-access-secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret"),
	}

	// Validate required fields
	if config.JWTAccessSecret == "your-access-secret" || config.JWTRefreshSecret == "your-refresh-secret" {
		log.Println("WARNING: Using default JWT secrets. Please set JWT_ACCESS_SECRET and JWT_REFRESH_SECRET in production.")
	}

	return config
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
