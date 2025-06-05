package config

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	TelegramBotToken     string
	TelegramChannelID    int64
	DatabaseDSN          string
	FetchInterval        time.Duration
	NotificationInterval time.Duration
	LookupTimeWindow     time.Duration
	OpenAIKey            string
	OpenAIPrompt         string
	OpenAIModel          string
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {

		// Load from .env file
		// err := godotenv.Load("../.env")
		// if err != nil {
		// 	log.Println("[WARN] .env file not found, falling back to system env vars")
		// }

		// Or load from docker

		// env is always string
		channelID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHANNEL_ID"), 10, 64)
		fetchInterval, _ := time.ParseDuration(os.Getenv("FETCH_INTERVAL"))
		notifyInterval, _ := time.ParseDuration(os.Getenv("NOTIFICATION_INTERVAL"))
		lookupTimeWindow, _ := time.ParseDuration(os.Getenv("LOOK_UP_TIME_WINDOW"))

		cfg = &Config{
			TelegramBotToken:     mustGet("TELEGRAM_BOT_TOKEN"),
			TelegramChannelID:    channelID,
			DatabaseDSN:          os.Getenv("DATABASE_DSN"),
			FetchInterval:        fetchInterval,
			NotificationInterval: notifyInterval,
			LookupTimeWindow:     lookupTimeWindow,
			OpenAIKey:            mustGet("OPENAI_KEY"),
			OpenAIPrompt:         os.Getenv("OPENAI_PROMPT"),
			OpenAIModel:          os.Getenv("OPENAI_MODEL"),
		}
	})

	return cfg
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return val
}
