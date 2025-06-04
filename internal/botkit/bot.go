package botkit

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
}

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api:      api,
		cmdViews: make(map[string]ViewFunc),
	}
}

func (b *Bot) RegisterCommand(cmd string, view ViewFunc) {
	b.cmdViews[cmd] = view
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	defer func() {
		if r := recover(); r != nil {
			// Handle panic, log it or notify
			log.Printf("[ERROR] recovered from panic: %v", r)
		}
	}()

	if update.Message == nil || !update.Message.IsCommand() {
		return nil
	}

	var view ViewFunc

	if !update.Message.IsCommand() {
		return nil // Ignore non-command messages
	}

	cmd := update.Message.Command()

	cmdView, exists := b.cmdViews[cmd]
	if !exists {
		return nil
	}

	view = cmdView

	if err := view(ctx, b.api, update); err != nil {
		// Handle error, log it or notify
		log.Printf("[ERROR] handling command %s: %v", cmd, err)

		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "internal error.")); err != nil {
			log.Printf("[ERROR] sending error message: %v", err)
		}

	}

	return nil
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Minute)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
