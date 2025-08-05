# 🤖 Redhat Bot - Telegram AI Assistant

A fully functional Telegram bot built in Go with AI capabilities, featuring daily group summaries, joke generation, and intelligent question answering.

## ✨ Features

### 🤖 AI Commands
- **`/redhat <question>`** - Ask the AI assistant anything
- **`/RTJ <topic>`** - Generate topic-based jokes
- **`/RRS`** - Check remaining daily usage
- **`/start`** - Welcome message and introduction
- **`/help`** - Detailed help and usage information

### 🧠 Daily Group Analysis
- Automatically logs all group messages (24h retention)
- Generates AI-powered daily summaries
- Posts summaries back to groups at 9 AM daily
- Filters out commands and short messages

### 🛡️ Rate Limiting & Security
- 5 requests per user per day
- 10-second cooldown between requests
- Automatic daily reset
- In-memory storage for performance

## 🏗️ Architecture

```
┌──────────────────────────────┐
│     Telegram Bot API         │
└────────────┬─────────────────┘
             │
             ▼
┌──────────────────────────────┐
│       Go Bot (tgbotapi)       │
├────────────┬─────────────────┤
│ Command Router                │
│ Group Message Listener        │
│ Daily Summarizer Scheduler    │
│ Rate Limiter Middleware       │
└───────┬────────────┬──────────┘
        │            │
        ▼            ▼
┌────────────────────┐    ┌────────────────────┐
│ In-Memory Storage   │    │ DeepSeek AI API    │
└────────────────────┘    └────────────────────┘
```

## 📦 Project Structure

```
redhat-bot/
├── main.go                 # Main bot entry point
├── config/
│   └── env.go             # Configuration management
├── commands/
│   ├── redhat.go          # /redhat command handler
│   ├── rtj.go             # /RTJ command handler
│   └── rrs.go             # /RRS command handler
├── limiter/
│   └── rate_limiter.go    # Rate limiting logic
├── storage/
│   └── memory.go          # In-memory data storage
├── ai/
│   └── openai.go          # DeepSeek AI integration
├── scheduler/
│   └── daily_summary.go   # Daily summary scheduler
├── go.mod                 # Go dependencies
└── README.md              # This file
```

## 🚀 Quick Start

### Prerequisites
- Go 1.24.5 or higher
- Telegram Bot Token
- DeepSeek API Token

### Installation

1. **Clone and navigate to the project:**
   ```bash
   cd redhat-bot
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables (optional):**
   Create a `.env` file:
   ```env
   TELEGRAM_TOKEN=your_telegram_bot_token
   DEEPSEEK_TOKEN=your_deepseek_api_token
   MAX_REQUESTS_PER_DAY=5
   COOLDOWN_SECONDS=10
   ```

4. **Run the bot:**
   ```bash
   go run main.go
   ```

### Default Configuration
The bot comes pre-configured with:
- **Telegram Token**: `8274702080:AAF7TqmXwUwTqWNCf1i52A5kLEG3uAQdqqE`
- **DeepSeek Token**: `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`
- **Daily Limit**: 5 requests per user
- **Cooldown**: 10 seconds between requests

## 🎯 Usage Examples

### Private Chat
```
User: /redhat What is the capital of France?
Bot: 🤖 Redhat AI

The capital of France is Paris. It's a beautiful city known for its rich history, culture, and iconic landmarks like the Eiffel Tower.

💡 Usage: 4/5 requests today
```

### Group Chat
```
User: /RTJ programming
Bot: 😄 Redhat Joke Generator

Topic: programming

Why do programmers prefer dark mode? Because light attracts bugs! 🐛

💡 Usage: 3/5 requests today
```

### Daily Summary (Automatic)
```
Bot: 🧠 Daily Group Summary

📅 Date: January 15, 2024
👥 Group: My Awesome Group
💬 Messages Analyzed: 47

Today the group discussed various programming topics, with Alice sharing tips about Go development and Bob asking about Docker containers. The conversation was lively and educational, ending with some helpful debugging advice.

Powered by Redhat Bot 🤖
```

## 🔧 Configuration

### Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `TELEGRAM_TOKEN` | Pre-configured | Your Telegram bot token |
| `DEEPSEEK_TOKEN` | Pre-configured | Your DeepSeek API token |
| `MAX_REQUESTS_PER_DAY` | 5 | Daily request limit per user |
| `COOLDOWN_SECONDS` | 10 | Cooldown between requests |

### Customization
- **Daily Summary Time**: Modify the cron schedule in `main.go` (currently 9 AM)
- **Message Retention**: Adjust the 24-hour window in `storage/memory.go`
- **AI Prompts**: Customize AI behavior in `ai/openai.go`

## 🛠️ Development

### Adding New Commands
1. Create a new command file in `commands/`
2. Implement the command handler
3. Add the command to the switch statement in `main.go`

### Extending AI Features
- Modify prompts in `ai/openai.go`
- Add new AI functions for different use cases
- Integrate additional AI providers

### Database Integration
The bot currently uses in-memory storage. To add persistence:
1. Implement database interfaces in `storage/`
2. Add database connection logic
3. Update storage methods to use the database

## 📊 Monitoring

The bot logs important events:
- Bot startup and configuration
- Command usage and rate limiting
- Daily summary generation
- Error messages and API failures

## 🔒 Security Features

- **Rate Limiting**: Prevents abuse with daily limits and cooldowns
- **Input Validation**: Sanitizes user inputs
- **Error Handling**: Graceful error handling without exposing sensitive data
- **Memory Management**: Automatic cleanup of old messages

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📝 License

This project is open source and available under the MIT License.

## 🆘 Support

If you encounter any issues:
1. Check the logs for error messages
2. Verify your API tokens are correct
3. Ensure the bot has proper permissions in groups
4. Check rate limiting if commands aren't working

---

**Built with ❤️ using Go and DeepSeek AI** 