package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
)

func main() {
	log.Println("DEPLOY MARK:", time.Now().UnixNano())

	token := os.Getenv("TELEGRAM_TOKEN")
	dsn := os.Getenv("DATABASE_URL")
	fmt.Printf("TELEGRAM_TOKEN: %s, DATABASE_URL: %s\n", token, dsn)

	// init database
	db := initDB(dsn)

	// health check handler
	go healthCheckAPIHandler(db)

	// use cases
	container := usecase.NewContainer(db)

	// init telegram app
	app, err := telegram.New(token, container)
	if err != nil {
		return
	}

	go initServer(container)

	// run telegram app
	if err = app.Run(); err != nil {
		log.Panic(err)
	}
}
