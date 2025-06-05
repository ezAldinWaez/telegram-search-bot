package bot

import (
	"fmt"
	"log"
	"telegram-semantic-search/config"
	"telegram-semantic-search/database"
	"telegram-semantic-search/embedding"
	"telegram-semantic-search/search"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api       *tgbotapi.BotAPI
	db        *database.DB
	config    *config.Config
	embedding *embedding.Client
	search    *search.Engine
}

func NewBot(cfg *config.Config, db *database.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = false
	log.Printf("Authorized on account %s", api.Self.UserName)

	// Initialize embedding client
	embeddingClient := embedding.NewClient(cfg.EmbeddingAPIURL, cfg.EmbeddingModel)

	// Initialize search engine
	searchEngine := search.NewEngine(db, embeddingClient, cfg.MaxResults)

	// Test embedding connection (non-blocking)
	go func() {
		if err := embeddingClient.TestConnection(); err != nil {
			log.Printf("⚠️  Embedding service connection failed: %v", err)
			log.Printf("💡 Make sure Ollama is running: ollama serve")
			log.Printf("💡 And model is available: ollama pull %s", cfg.EmbeddingModel)
		} else {
			log.Printf("✅ Embedding service connected successfully")
		}
	}()

	return &Bot{
		api:       api,
		db:        db,
		config:    cfg,
		embedding: embeddingClient,
		search:    searchEngine,
	}, nil
}

func (b *Bot) Start() error {
	log.Println("Starting bot...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		go b.handleUpdate(update)
	}

	return nil
}

func (b *Bot) Stop() {
	log.Println("Stopping bot...")
	b.api.StopReceivingUpdates()
}
