package bot

import (
	"context"

	"github.com/amir-amirov/go-news-feed-bot/internal/botkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Welcome to the News Feed Bot!")); err != nil {
			return err
		}
		return nil
	}
}
