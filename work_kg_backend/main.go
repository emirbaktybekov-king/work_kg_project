package main

import (
	"log"

	"work_kg_backend/internal/bot"
	"work_kg_backend/internal/config"
	"work_kg_backend/internal/database"
	"work_kg_backend/internal/handlers"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize database schema
	database.InitSchema()

	// Start Telegram bot in goroutine
	go bot.Start(cfg.TelegramToken)

	// Start HTTP server (blocking)
	handlers.StartServer(cfg.ServerPort)
}
