# Telegram Semantic Search Bot - Phase 1

A Telegram bot that will eventually provide semantic search capabilities over chat messages. Currently in Phase 1: Basic bot + database functionality.

## Features (Phase 1)

-   ✅ Connect to Telegram and receive messages
-   ✅ Store messages in SQLite database
-   ✅ Basic text preprocessing and cleaning
-   ✅ Bot commands: `/start`, `/help`, `/stats`
-   ✅ Message tracking and storage
-   🔄 Embedding generation (Phase 2)
-   🔄 Semantic search (Phase 3)

## Setup

### 1. Prerequisites

-   Go 1.21 or higher
-   SQLite3 (included with go-sqlite3 driver)

### 2. Get Telegram Bot Token

1. Message [@BotFather](https://t.me/botfather) on Telegram
2. Create a new bot with `/newbot`
3. Copy the bot token

### 3. Installation

```bash
# Clone or create the project
mkdir telegram-semantic-search
cd telegram-semantic-search

# Initialize Go module
go mod init telegram-semantic-search

# Install dependencies and setup
make setup

# Edit .env file with your bot token
# TELEGRAM_TOKEN=your_actual_bot_token_here

# Run the bot
make run
```

### 4. Configuration

The bot uses a `.env` file for configuration:

```bash
# Copy the example file (done automatically by 'make setup')
cp .env.example .env

# Edit .env file with your settings
TELEGRAM_TOKEN=your_bot_token_here
DATABASE_PATH=./messages.db
EMBEDDING_API_URL=http://localhost:11434
EMBEDDING_MODEL=all-minilm:latest
```

## Usage

### 1. Add Bot to Chat

1. Find your bot by username (given by BotFather)
2. Add it to a group or start a private chat
3. Send `/start` to initialize

### 2. Available Commands

-   `/start` - Welcome message and setup
-   `/help` - Show help information
-   `/stats` - Show message statistics for current chat
-   `/search <query>` - Placeholder (Phase 3)

### 3. Message Tracking

The bot will automatically:

-   Track all text messages in chats where it's added
-   Store messages with metadata (user, timestamp, chat)
-   Clean and preprocess text
-   Log activity (check console output)

## Database Schema

```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    username TEXT,
    text TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    embedding TEXT -- For Phase 2
);
```

## Testing Phase 1

1. **Message Storage**: Send various messages and check `/stats`
2. **Commands**: Test all available commands
3. **Multiple Chats**: Add bot to different groups/chats
4. **Database**: Check `messages.db` file is created and growing

## What's Next

-   **Phase 2**: Add embedding generation with Ollama
-   **Phase 3**: Implement semantic search functionality
-   **Phase 4**: Testing and refinement

## Troubleshooting

### Common Issues

1. **Bot not responding**: Check TELEGRAM_TOKEN is correct
2. **Database errors**: Ensure write permissions in current directory
3. **No messages stored**: Bot needs to be added to chat and see new messages

### Logs

The bot logs important events to console:

-   Message storage confirmations
-   Command executions
-   Database operations
-   Errors and warnings

## Architecture

```
main.go
├── config/ - Configuration management
├── database/ - SQLite operations and models
└── bot/ - Telegram bot logic and handlers
```

This Phase 1 implementation provides a solid foundation for adding semantic search capabilities in the next phases.
