## Architecture Overview

```
Telegram Message → Embedding Generation → Vector Storage → Semantic Search
                ↓                        ↓               ↓
           Text Processing          SQLite Database    Cosine Similarity
                ↓                        ↓               ↓
           Ollama API              JSON Embeddings    Ranked Results
```

### Core Components

-   **Bot Layer**: Telegram API integration and command handling
-   **Embedding Layer**: Ollama integration for vector generation
-   **Database Layer**: SQLite storage with JSON embedding fields
-   **Search Engine**: Vector similarity calculations and ranking
-   **Configuration**: Environment-based setup with .env support

## Success Metrics

**Primary Metric: Search Relevance Rate**

-   Target: >70% of searches return at least one relevant message in top 3
-   Current: Ready for testing with real conversations

**Key Performance Indicators:**

-   Search response time: <2 seconds for typical chat history
-   Embedding generation: Background processing, doesn't block chat
-   Storage efficiency: JSON embeddings in SQLite, ~384 dims per message
-   Memory usage: Optimized for moderate chat volumes (1000s of messages)

## Troubleshooting

### Common Issues

1. **Search returns no results**

    - Check if messages have embeddings: `/stats`
    - Try broader search terms
    - Ensure minimum 3+ character messages are being processed

2. **Slow search performance**

    - Ollama service might be slow
    - Large chat history (>10k messages) may need optimization
    - Check embedding generation backlog in logs

3. **Poor search relevance**
    - Need more diverse messages for better training
    - Try different phrasing of search query
    - Some topics may not have been discussed yet

### Performance Optimization

-   **Database**: Automatic indexing on chat_id and timestamp
-   **Memory**: Messages loaded only when searching, not kept in memory
-   **Concurrency**: Embedding generation runs in background goroutines
-   **Caching**: Consider adding query result caching for frequently searched terms

## Development Notes

The MVP successfully implements all core requirements:

-   ✅ Message tracking and storage
-   ✅ Semantic embedding generation
-   ✅ Vector similarity search
-   ✅ Telegram bot interface
-   ✅ Configurable and maintainable codebase

**Total codebase**: ~800 lines of clean, well-structured Go code
**Dependencies**: Minimal and stable (tgbotapi, sqlite3, godotenv)
**Deployment**: Single binary, self-contained with SQLite# Telegram Semantic Search Bot - Phase 3 ✅

A Telegram bot that provides semantic search capabilities over chat messages. **Phase 3 Complete: Fully operational semantic search!**

## Features (Phase 3)

-   ✅ Connect to Telegram and receive messages
-   ✅ Store messages in SQLite database
-   ✅ Basic text preprocessing and cleaning
-   ✅ Bot commands: `/start`, `/help`, `/stats`, `/test`, **`/search`**
-   ✅ Message tracking and storage
-   ✅ Semantic embedding generation using Ollama
-   ✅ Asynchronous embedding processing
-   ✅ Embedding service health checks
-   ✅ **Semantic search with cosine similarity**
-   ✅ **Intelligent result ranking and formatting**
-   ✅ **Natural language query understanding**

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
-   `/help` - Show help information and search examples
-   `/stats` - Show message and embedding statistics
-   `/test` - Test embedding service connection
-   **`/search <query>` - Semantic search through chat history**

### 3. Semantic Search Examples

```
/search meeting tomorrow     # Finds discussions about meetings
/search python programming   # Finds code-related conversations
/search funny joke          # Finds humorous messages
/search project deadline     # Finds work planning discussions
/search weekend plans        # Finds casual planning talks
```

**Key Features:**

-   **Semantic understanding** - searches by meaning, not just keywords
-   **Ranked results** - shows similarity percentages
-   **Context-aware** - understands related concepts
-   **Fast and efficient** - optimized vector similarity search

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

## Testing Phase 3

1. **Setup Ollama**: Ensure `ollama serve` is running and `all-minilm` model is pulled
2. **Add Bot to Chat**: Invite bot to a group or start private chat
3. **Generate Content**: Send various messages to build search index
4. **Test Search**: Try `/search <query>` with different topics
5. **Check Results**: Verify relevance and similarity scores

### Search Testing Examples

```bash
# Add bot to chat and send these messages:
"Let's schedule a team meeting for tomorrow at 3 PM"
"I'm working on a Python script for data analysis"
"That joke was hilarious! 😂"
"The project deadline is next Friday"
"What are your plans for the weekend?"

# Then test searches:
/search meeting schedule     # Should find the meeting message
/search python code         # Should find the programming message
/search funny               # Should find the joke message
/search deadline            # Should find the project message
/search weekend             # Should find the weekend plans message
```

### Search Quality Indicators

-   **High similarity (>70%)**: Very relevant results
-   **Medium similarity (30-70%)**: Somewhat related results
-   **Low similarity (<30%)**: Filtered out automatically
-   **No results**: Need more messages or different search terms

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
