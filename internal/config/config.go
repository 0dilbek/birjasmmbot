package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken      string
	DatabaseURL   string
	AdminIDs      []int64
	FreeResponses int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		BotToken:      mustGetenv("BOT_TOKEN"),
		DatabaseURL:   mustGetenv("DATABASE_URL"),
		FreeResponses: getenvInt("FREE_RESPONSES", 5),
	}

	for _, idStr := range strings.Split(os.Getenv("ADMIN_IDS"), ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			cfg.AdminIDs = append(cfg.AdminIDs, id)
		}
	}

	return cfg
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Required env variable %s is not set", key)
	}
	return v
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
