package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"telegram-semantic-search/bot"
	"telegram-semantic-search/config"
	"telegram-semantic-search/database"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN environment variable is required")
	}

	// Initialize database
	db, err := database.NewDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Printf("Database initialized at: %s", cfg.DatabasePath)

	// Initialize bot
	telegramBot, err := bot.NewBot(cfg, db)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received shutdown signal")
		telegramBot.Stop()
		os.Exit(0)
	}()

	// Start bot
	log.Println("Starting Telegram Semantic Search Bot...")
	log.Println("Phase 3: Semantic search fully operational! 🎯")
	log.Printf("Embedding model: %s", cfg.EmbeddingModel)
	log.Printf("Embedding API: %s", cfg.EmbeddingAPIURL)
	log.Printf("Max search results: %d", cfg.MaxResults)
	log.Println("Use /start command to interact with the bot")
	log.Println("Use /search <query> to perform semantic search!")

	if err := telegramBot.Start(); err != nil {
		log.Fatalf("Bot failed: %v", err)
	}
}
