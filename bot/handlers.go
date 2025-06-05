package bot

import (
	"fmt"
	"log"
	"strings"
	"telegram-semantic-search/database"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	// Handle regular messages
	if update.Message != nil {
		b.handleMessage(update.Message)
		return
	}

	// Handle edited messages (optional - treat as new messages)
	if update.EditedMessage != nil {
		b.handleMessage(update.EditedMessage)
		return
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Skip empty messages
	if message.Text == "" {
		return
	}

	// Handle commands
	if message.IsCommand() {
		b.handleCommand(message)
		return
	}

	// Store regular messages
	b.storeMessage(message)
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	command := message.Command()

	switch command {
	case "start":
		b.handleStartCommand(message)
	case "help":
		b.handleHelpCommand(message)
	case "stats":
		b.handleStatsCommand(message)
	case "test":
		b.handleTestCommand(message)
	case "search":
		// Placeholder for Phase 3
		b.sendReply(message, "🔍 Search functionality will be added in Phase 3!")
	default:
		b.sendReply(message, fmt.Sprintf("Unknown command: /%s", command))
	}

	log.Printf("Command /%s executed by %s in chat %d", command, message.From.UserName, message.Chat.ID)
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) {
	welcomeText := `🤖 *Semantic Search Bot*

I'm now tracking messages in this chat and generating semantic embeddings for each one!

*Available commands:*
• /help - Show help message
• /stats - Show message and embedding statistics
• /test - Test embedding service connection
• /search <query> - Search messages (Phase 3)

*Phase 2 Active:* I'm now generating embeddings for all messages using AI. This enables semantic search capabilities!

Just keep chatting normally - I'll process and index your messages automatically! 🚀`

	b.sendReply(message, welcomeText)
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `🤖 *Semantic Search Bot Help*

*What I do:*
• Track all messages in this chat
• Generate semantic embeddings for each message
• Help you find relevant past conversations (Phase 3)

*Commands:*
• /start - Welcome message
• /help - This help message  
• /stats - Show message and embedding statistics
• /test - Test embedding service connection
• /search <query> - Semantic search (Phase 3)

*Current Status:* Phase 2 - Embedding generation active
*Privacy:* Messages stored locally, used only for search functionality.`

	b.sendReply(message, helpText)
}

func (b *Bot) handleStatsCommand(message *tgbotapi.Message) {
	count, err := b.db.GetStats(message.Chat.ID)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		b.sendReply(message, "❌ Error getting statistics")
		return
	}

	// Count messages with embeddings
	countWithEmbeddings, err := b.db.GetStatsWithEmbeddings(message.Chat.ID)
	if err != nil {
		log.Printf("Error getting embedding stats: %v", err)
		countWithEmbeddings = 0
	}

	statsText := fmt.Sprintf(`📊 *Chat Statistics*

Messages stored: *%d*
Messages with embeddings: *%d*
Chat ID: %d
Embedding model: %s
Status: ✅ Active

Ready for semantic search once Phase 3 is complete!`, count, countWithEmbeddings, message.Chat.ID, b.config.EmbeddingModel)

	b.sendReply(message, statsText)
}

func (b *Bot) handleTestCommand(message *tgbotapi.Message) {
	b.sendReply(message, "🧪 Testing embedding service...")

	// Test embedding generation
	testText := "This is a test message for embedding generation"
	embedding, err := b.embedding.GetEmbedding(testText)
	if err != nil {
		errorMsg := fmt.Sprintf(`❌ *Embedding Test Failed*

Error: %s

*Troubleshooting:*
• Make sure Ollama is running: `+"`ollama serve`"+`
• Check if model is available: `+"`ollama pull %s`"+`
• Verify API URL: %s`, err.Error(), b.config.EmbeddingModel, b.config.EmbeddingAPIURL)

		b.sendReply(message, errorMsg)
		return
	}

	successMsg := fmt.Sprintf(`✅ *Embedding Test Successful*

Test text: "%s"
Embedding dimensions: *%d*
Model: %s
API URL: %s

Embedding service is working correctly!`, testText, len(embedding), b.config.EmbeddingModel, b.config.EmbeddingAPIURL)

	b.sendReply(message, successMsg)
}

func (b *Bot) storeMessage(message *tgbotapi.Message) {
	// Clean the message text (basic preprocessing)
	cleanText := b.cleanText(message.Text)

	// Skip very short messages
	if len(strings.TrimSpace(cleanText)) < 3 {
		return
	}

	// Create message object
	msg := database.Message{
		ChatID:    message.Chat.ID,
		UserID:    message.From.ID,
		Username:  message.From.UserName,
		Text:      cleanText,
		Timestamp: time.Unix(int64(message.Date), 0),
	}

	// Generate embedding asynchronously
	go func() {
		embedding, err := b.embedding.GetEmbedding(cleanText)
		if err != nil {
			log.Printf("Failed to generate embedding for message: %v", err)
			// Save message without embedding
			if err := b.db.SaveMessage(msg); err != nil {
				log.Printf("Error saving message without embedding: %v", err)
			}
			return
		}

		// Add embedding to message
		msg.Embedding = embedding

		// Save message with embedding
		if err := b.db.SaveMessage(msg); err != nil {
			log.Printf("Error saving message with embedding: %v", err)
		} else {
			log.Printf("✅ Saved message with embedding (%d dims) from %s", len(embedding), msg.Username)
		}
	}()
}

func (b *Bot) cleanText(text string) string {
	// Basic text cleaning
	// Remove multiple spaces
	text = strings.Join(strings.Fields(text), " ")

	// Trim whitespace
	text = strings.TrimSpace(text)

	return text
}

func (b *Bot) sendReply(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
