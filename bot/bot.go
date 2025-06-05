package bot

import (
	"fmt"
	"log"
	"telegram-semantic-search/config"
	"telegram-semantic-search/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	db     *database.DB
	config *config.Config
}

func NewBot(cfg *config.Config, db *database.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = false
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:    api,
		db:     db,
		config: cfg,
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
