package bot

import (
	"fmt"
	"log"
	"strings"
	"telegram-semantic-search/database"
	"telegram-semantic-search/search"
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
	args := message.CommandArguments()

	switch command {
	case "start":
		b.handleStartCommand(message)
	case "help":
		b.handleHelpCommand(message)
	case "stats":
		b.handleStatsCommand(message)
	case "test":
		b.handleTestCommand(message)
	case "perf":
		b.handlePerfCommand(message)
	case "search":
		b.handleSearchCommand(message, args)
	default:
		b.sendReply(message, fmt.Sprintf("Unknown command: /%s", command))
	}

	log.Printf("Command /%s executed by %s in chat %d", command, message.From.UserName, message.Chat.ID)
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) {
	welcomeText := `🤖 *Semantic Search Bot*

I'm now fully operational with semantic search capabilities! 🚀

*Available commands:*
• /help - Show help message
• /stats - Show message and embedding statistics
• /test - Test embedding service connection
• /perf - Show performance statistics  
• /search <query> - **Search messages semantically**

*Phase 4 Complete:* You can now search through chat history using natural language!

Examples:
• ` + "`/search meeting tomorrow`" + `
• ` + "`/search funny story`" + `
• ` + "`/search project deadline`" + `

Just keep chatting - I'll continue indexing messages for better search results! 🎯`

	b.sendReply(message, welcomeText)
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `🤖 *Semantic Search Bot Help*

*What I do:*
• Track all messages in this chat
• Generate semantic embeddings for each message
• **Enable semantic search through chat history**

*Commands:*
• /start - Welcome message
• /help - This help message  
• /stats - Show message and embedding statistics
• /test - Test embedding service connection
• /perf - Show performance statistics
• /search <query> - **Search messages by meaning, not just keywords**

*How search works:*
I understand context and meaning, not just exact word matches!

*Examples:*
• ` + "`/search meeting schedule`" + ` - finds discussions about meetings
• ` + "`/search python code`" + ` - finds programming conversations  
• ` + "`/search weekend plans`" + ` - finds casual planning discussions

*Current Status:* Phase 4 - **Production-ready semantic search**
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
Search readiness: *%.1f%%*
Chat ID: %d
Embedding model: %s
Status: ✅ **Semantic search active**

Ready to search! Try: `+"`/search your query`"+``, count, countWithEmbeddings, float64(countWithEmbeddings)/float64(max(count, 1))*100, message.Chat.ID, b.config.EmbeddingModel)

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

func (b *Bot) handlePerfCommand(message *tgbotapi.Message) {
	searchAvg, embeddingAvg, memUsage := b.perf.GetStats()

	perfMsg := fmt.Sprintf(`⚡ *Performance Statistics*

*Search Performance:*
• Average search time: %v
• Target: < 2 seconds

*Embedding Performance:*
• Average embedding time: %v
• Background processing

*System Resources:*
• Memory usage: %s
• Status: %s

*Optimization Notes:*
• Search performance scales with message count
• Embedding generation runs asynchronously
• Memory usage optimized for chat volumes`,
		searchAvg,
		embeddingAvg,
		memUsage,
		getPerformanceStatus(searchAvg))

	b.sendReply(message, perfMsg)
}

func (b *Bot) handleSearchCommand(message *tgbotapi.Message, query string) {
	if strings.TrimSpace(query) == "" {
		b.sendReply(message, `🔍 *Semantic Search*

Usage: `+"`/search <your query>`"+`

Examples:
• `+"`/search python programming`"+`
• `+"`/search meeting schedule`"+`
• `+"`/search funny joke`"+`

I'll find the most relevant messages based on semantic similarity!`)
		return
	}

	// Show searching indicator
	b.sendReply(message, fmt.Sprintf("🔍 Searching for: *%s*...", query))

	// Start performance timing
	startTime := time.Now()

	// Perform search
	results, err := b.search.Search(query, message.Chat.ID)

	// Record search performance
	searchDuration := time.Since(startTime)
	b.perf.RecordSearchTime(searchDuration)

	if err != nil {
		log.Printf("Search error: %v", err)
		b.sendReply(message, fmt.Sprintf("❌ Search failed: %s", err.Error()))
		return
	}

	// Handle no results
	if len(results) == 0 {
		totalMessages, withEmbeddings, _ := b.search.SearchStats(message.Chat.ID)

		noResultsMsg := fmt.Sprintf(`🤷‍♂️ *No Results Found*

Query: "%s"

*Possible reasons:*
• No similar messages found (similarity too low)
• Not enough messages with embeddings yet
• Try different search terms

*Chat Stats:*
• Total messages: %d
• Messages with embeddings: %d

Try searching for topics you know were discussed!`, query, totalMessages, withEmbeddings)

		b.sendReply(message, noResultsMsg)
		return
	}

	// Format and send results
	resultMsg := b.formatSearchResults(query, results, searchDuration)
	b.sendReply(message, resultMsg)

	log.Printf("Search completed: query='%s', results=%d, duration=%v, chat=%d",
		query, len(results), searchDuration, message.Chat.ID)
}

func (b *Bot) formatSearchResults(query string, results []search.SearchResult, searchDuration time.Duration) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("🎯 *Search Results for:* \"%s\"\n", query))
	msg.WriteString(fmt.Sprintf("⚡ *Search time:* %v\n\n", searchDuration))

	for _, result := range results {
		// Format timestamp
		timeStr := result.Message.Timestamp.Format("Jan 2, 15:04")

		// Truncate long messages
		text := result.Message.Text
		if len(text) > 200 {
			text = text[:200] + "..."
		}

		// Format similarity percentage
		similarity := fmt.Sprintf("%.1f%%", result.Similarity*100)

		msg.WriteString(fmt.Sprintf("**%d.** *%s* (%s similarity)\n",
			result.Rank, similarity, timeStr))
		msg.WriteString(fmt.Sprintf("👤 %s\n", result.Message.Username))
		msg.WriteString(fmt.Sprintf("💬 %s\n\n", text))
	}

	msg.WriteString("_💡 Results ranked by semantic similarity_")

	return msg.String()
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
		startTime := time.Now()

		embedding, err := b.embedding.GetEmbedding(cleanText)

		// Record embedding performance
		embeddingDuration := time.Since(startTime)
		b.perf.RecordEmbeddingTime(embeddingDuration)

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
			log.Printf("✅ Saved message with embedding (%d dims, %v) from %s",
				len(embedding), embeddingDuration, msg.Username)
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

func getPerformanceStatus(searchAvg time.Duration) string {
	if searchAvg == 0 {
		return "🟡 No searches performed yet"
	} else if searchAvg < 2*time.Second {
		return "🟢 Excellent"
	} else if searchAvg < 5*time.Second {
		return "🟡 Good"
	} else {
		return "🔴 Needs optimization"
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
