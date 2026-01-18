# Training Telegram Bot

## Overview
A Telegram bot for tracking workout training sessions. Built with Go and uses SQLite for data storage.

## Project Structure
- `cmd/main.go` - Main entry point
- `internal/config/` - Configuration handling
- `internal/models/` - Database models (User, WorkoutDay, Exercise, Set, WorkoutSession)
- `internal/repository/` - Database repositories for each model
- `internal/service/` - Bot service and message handlers
- `internal/templates/` - Message templates
- `internal/utils/` - Utility functions

## Dependencies
- Go 1.24
- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot API
- `gorm.io/gorm` with SQLite driver - ORM and database

## Configuration
The bot requires:
1. Environment variable `TELEGRAM_TOKEN` containing your Telegram Bot API token
2. Environment variable `DATABASE_DSN` containing DSN for connection to your database

## Running
```bash
go build -o training-tg-bot ./cmd/main.go
./training-tg-bot
```

## Database
Uses Postgres database. Auto-migrates on startup via github.com/pressly/goose/v3.

## Required Secrets
- `TELEGRAM_TOKEN`: Your Telegram Bot API token (get from @BotFather on Telegram)

## PG_DSN for local startup
postgresql://postgres:postgres@localhost/training-bot
