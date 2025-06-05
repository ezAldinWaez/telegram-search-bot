# Telegram Semantic Search Bot - Phase 2

A Telegram bot that provides semantic search capabilities over chat messages. Currently in Phase 2: Embedding generation with Ollama.

## Features (Phase 2)

-   ✅ Connect to Telegram and receive messages
-   ✅ Store messages in SQLite database
-   ✅ Basic text preprocessing and cleaning
-   ✅ Bot commands: `/start`, `/help`, `/stats`, `/test`
-   ✅ Message tracking and storage
-   ✅ **Semantic embedding generation using Ollama**
-   ✅ **Asynchronous embedding processing**
-   ✅ **Embedding service health checks**
-   🔄 Semantic search (Phase 3)

## Setup

### 1. Prerequisites

-   Go 1.21 or higher
-   SQLite3 (included with go-sqlite3 driver)
-   **Ollama** (for embedding generation)

### 2. Setup Ollama

```bash
# Install Ollama (if not already installed)
# Visit: https://ollama.ai/download

# Start Ollama service
ollama serve

# Pull the embedding model (in another terminal)
ollama pull all-minilm

# Verify model is available
ollama list
```

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
-   `/stats` - Show message and embedding statistics
-   `/test` - Test embedding service connection
-   `/search <query>` - Placeholder (Phase 3)

### 3. Message Processing

The bot will automatically:

-   Track all text messages in chats where it's added
-   Store messages with metadata (user, timestamp, chat)
-   **Generate semantic embeddings for each message**
-   **Process embeddings asynchronously (non-blocking)**
-   Clean and preprocess text
-   Log activity and embedding status

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

## Testing Phase 2

1. **Setup Ollama**: Ensure `ollama serve` is running and `all-minilm` model is pulled
2. **Test Embedding Service**: Use `/test` command to verify connection
3. **Message Processing**: Send messages and check console for embedding confirmations
4. **Statistics**: Use `/stats` to see embedding generation progress
5. **Multiple Chats**: Test embedding generation across different groups/chats

### Troubleshooting Embeddings

-   **"Connection failed"**: Make sure `ollama serve` is running
-   **"Model not found"**: Run `ollama pull all-minilm`
-   **Slow processing**: Embeddings are generated asynchronously - check logs
-   **Missing embeddings**: Some messages may be saved without embeddings if API fails

## What's Next

-   **Phase 3**: Implement semantic search functionality
-   **Phase 4**: Testing and refinement

## New in Phase 2

-   **Ollama Integration**: Automatic embedding generation for all messages
-   **Async Processing**: Non-blocking embedding generation
-   **Health Checks**: Connection testing and model verification
-   **Enhanced Stats**: Track messages with/without embeddings
-   **Error Handling**: Graceful fallback when embedding service unavailable

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
