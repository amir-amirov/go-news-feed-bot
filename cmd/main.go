package main

import (
	"context"
	"log"

	"github.com/amir-amirov/go-news-feed-bot/internal/bot"
	"github.com/amir-amirov/go-news-feed-bot/internal/botkit"
	"github.com/amir-amirov/go-news-feed-bot/internal/config"
	"github.com/amir-amirov/go-news-feed-bot/internal/db"
	"github.com/amir-amirov/go-news-feed-bot/internal/fetcher"
	"github.com/amir-amirov/go-news-feed-bot/internal/notifier"
	"github.com/amir-amirov/go-news-feed-bot/internal/storage"
	"github.com/amir-amirov/go-news-feed-bot/internal/summary"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	config := config.Load()

	botAPI, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Fatalf("Failed to create bot API: %v", err)
	}
	log.Printf("Bot authorized as %s", botAPI.Self.UserName)

	db.InitDB(config.DatabaseDSN)

	filterKeywords := []string{"leetcode"}

	var (
		sourceRespository  = storage.NewSourcePostgresStorage(db.DB)
		articleRespository = storage.NewArticlePostgresStorage(db.DB)
		fetcher            = fetcher.New(articleRespository, sourceRespository, config.FetchInterval, filterKeywords)

		summarizer = summary.NewOpenAISummarizer(
			config.OpenAIKey,
			// "",
			config.OpenAIModel,
			config.OpenAIPrompt,
		)

		notifier = notifier.New(
			articleRespository,
			summarizer,
			botAPI,
			config.NotificationInterval,
			config.LookupTimeWindow,
			config.TelegramChannelID,
		)
	)

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCommand("start", bot.ViewCmdStart())

	go func() {
		fetcher.Start(context.TODO())
	}()

	go func() {
		notifier.Start(context.TODO())
	}()

	if err := newsBot.Run(context.TODO()); err != nil {
		log.Printf("[ERROR] failed to run botkit: %v", err)
	}
}
