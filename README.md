# Discord Quote Bot

* this is ai slop we love it *

A Discord bot that allows users to quote others with a slash command. The quotes are posted in a designated quotes channel with an embed containing the quote and a ping to the quoted user.

## Features

- `/quote` slash command to quote a user
- Modular command structure
- YAML configuration
- Clean and structured codebase

## Requirements

- Go 1.20 or higher
- A Discord Bot Token

## Setup

1. Clone the repository:
```bash
git clone https://github.com/txj-xyz/discord-quote-bot.git
cd discord-quote-bot
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure the bot:

Edit the `config.yaml` file:
```yaml
bot:
  token: "YOUR_BOT_TOKEN"
  # Uncomment for guild-specific testing
  # guild_id: "YOUR_GUILD_ID"

quote:
  channel: "YOUR_QUOTES_CHANNEL_ID"
```

4. Build and run the bot:
```bash
go build -o quote-bot .
./quote-bot
```

Or run directly:
```bash
go run main.go
```

## Usage

Once the bot is running and has joined your Discord server, you can start issuing commands.

## Project Structure

```
discord-quote-bot/
├── config/             # Configuration package
├── internal/           # Internal application code
│   ├── bot/            # Bot logic
│   └── commands/       # Command implementations
│   main.go             # Main bot entry point
└── config.yaml         # Bot configuration
``` 