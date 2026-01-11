package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/pressly/goose/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	_ "github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service"
)

func main() {
	fmt.Println("TELEGRAM_TOKEN:", os.Getenv("TELEGRAM_TOKEN"))
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("dsn:", dsn)

	if err = migrate(dsn); err != nil {
		log.Panicf("Failed to migrate database: %s", err.Error())
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect database: %s", err.Error())
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	usersRepo := users.NewRepo(db)
	programsRepo := programs.NewRepo(db)
	dayTypesRepo := daytypes.NewRepo(db)
	workoutsRepo := workouts.NewRepo(db)
	exercisesRepo := exercises.NewRepo(db)
	setsRepo := sets.NewRepo(db)
	sessionsRepo := sessions.NewRepo(db)
	exerciseTypesRepo := exercisetypes.NewRepo(db)
	exerciseGroupTypesRepo := exercisegrouptypes.NewRepo(db)

	svc := service.NewService(bot, usersRepo, programsRepo, dayTypesRepo, workoutsRepo, exercisesRepo,
		exerciseTypesRepo, exerciseGroupTypesRepo, setsRepo, sessionsRepo)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	fmt.Println("port:", port)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		})
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
			log.Printf("Health check server failed: %s", err)
		}
	}()

	for update := range updates {
		if update.Message != nil {
			svc.HandleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			svc.HandleCallback(update.CallbackQuery)
		}
	}
}

func migrate(dsn string) error {
	migrateDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer migrateDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(migrateDB, "db/migrations"); err != nil {
		return err
	}

	return nil
}
