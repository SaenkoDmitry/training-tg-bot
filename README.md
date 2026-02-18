# Training Telegram Bot

## Overview
A Telegram bot for tracking workout training sessions.

## Interface

### Preview and workout session
<p align="left">
  <img width="200" height="450" src="/screenshots/preview.png">  
  <img width="200" height="450" src="/screenshots/list_of_workouts+statistics.png">  
  <img width="200" height="450" src="/screenshots/workout_session.png">
</p>

### Export to Excel
<p align="left">
  <img width="300" height="150" src="/screenshots/excel_0.png">
  <img width="300" height="150" src="/screenshots/excel_1.png">
</p>

## Project Structure
- `cmd/main.go` - Main entry point
- `internal/models/` - Database models (User, WorkoutDay, Exercise, Set, WorkoutSession)
- `internal/repository/` - Database repositories for each model
- `internal/service/` - Bot service and message handlers
- `internal/utils/` - Utility functions

## Dependencies
- Go 1.24
- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot API
- `github.com/pressly/goose/v3` - for database migrations
- `gorm.io/gorm, gorm.io/driver/postgres` - ORM with Postgres driver

## Configuration and secrets
The bot requires:
1. Environment variable `TELEGRAM_TOKEN` containing your Telegram Bot API token (get from @BotFather on Telegram)
2. Environment variable `DATABASE_URL` containing DSN for connection to your database, for example 'postgresql://postgres:postgres@127.0.0.1/training_app_db?sslmode\=disable'
3. Environment variable `JWT_SECRET` containing JWT secret key
4. Environment variable `TELEGRAM_BOT_ID` containing your telegram bot id
5. Environment variable `VAPID_PRIVATE_KEY` and Environment variable `VAPID_PUBLIC_KEY` containing public/private keys for push notifications

## Running
```bash
go build -o training-tg-bot ./cmd/main.go
./training-tg-bot
```

## Database
Uses Postgres database. Auto-migrates on startup via github.com/pressly/goose/v3.

## DATABASE_URL for local startup
postgresql://postgres:postgres@localhost/training-bot

## how to local login to telegram?

1. configure tunnel
```
ssh -R 80:localhost:5173 localhost.run
```

2. edit domain in telegram for bot and paste host to it

