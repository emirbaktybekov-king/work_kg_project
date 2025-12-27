package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DatabaseURL   string
	ServerPort    string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		TelegramToken: getEnv("TELEGRAM_TOKEN", ""),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
	}

	// Validate required fields
	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
