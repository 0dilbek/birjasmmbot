package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/birjasmm/bot/internal/bot"
	"github.com/birjasmm/bot/internal/config"
	"github.com/birjasmm/bot/internal/db"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	b, err := bot.New(cfg, database)
	if err != nil {
		log.Fatalf("Bot init error: %v", err)
	}

	log.Println("Bot started. Press Ctrl+C to stop.")
	go b.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}
