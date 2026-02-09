package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	_ "github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB(dsn string) *gorm.DB {
	migrate(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect database: %s", err.Error())
	}
	return db
}

func healthCheckAPIHandler(db *gorm.DB) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	dbObj, err := db.DB()
	if err != nil {
		log.Fatalf("cannot get *sql.DB")
	}
	fmt.Println("port:", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err = dbObj.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB ping failed")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Printf("Health check server failed: %s", err)
	}
}

func migrate(dsn string) {
	migrateDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panicf("Failed to migrate: open dsn: %s", err.Error())
	}
	defer migrateDB.Close()

	if err = goose.SetDialect("postgres"); err != nil {
		log.Panicf("Failed to migrate: set database dialect: %s", err.Error())
	}

	if err := goose.Up(migrateDB, "db/migrations"); err != nil {
		log.Panicf("Failed to migrate: goose up: %s", err.Error())
	}
}
