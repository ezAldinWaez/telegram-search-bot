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

I'm now tracking messages in this chat and will soon be able to help you search through them semantically!

*Available commands:*
• /help - Show this help message
• /stats - Show message statistics
• /search <query> - Search messages (coming in Phase 3)

Just keep chatting normally - I'll silently collect and index your messages for semantic search!`

	b.sendReply(message, welcomeText)
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `🤖 *Semantic Search Bot Help*

*What I do:*
• Track all messages in this chat
• Store them for semantic search (coming soon)
• Help you find relevant past conversations

*Commands:*
• /start - Welcome message
• /help - This help message  
• /stats - Show how many messages I've stored
• /search <query> - Semantic search (Phase 3)

*Privacy:* I only store messages from chats where I'm added. Messages are stored locally and used only for search functionality.`

	b.sendReply(message, helpText)
}

func (b *Bot) handleStatsCommand(message *tgbotapi.Message) {
	count, err := b.db.GetStats(message.Chat.ID)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		b.sendReply(message, "❌ Error getting statistics")
		return
	}

	statsText := fmt.Sprintf(`📊 *Chat Statistics*

Messages stored: *%d*
Chat ID: %d
Status: ✅ Active

Ready for semantic search once Phase 3 is complete!`, count, message.Chat.ID)

	b.sendReply(message, statsText)
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
		// Embedding will be added in Phase 2
	}

	// Save to database
	if err := b.db.SaveMessage(msg); err != nil {
		log.Printf("Error saving message: %v", err)
	}
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
