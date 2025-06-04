package notifier

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/amir-amirov/go-news-feed-bot/internal/botkit/markup"
	"github.com/amir-amirov/go-news-feed-bot/internal/model"
	"github.com/go-shiori/go-readability"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArticlesRepository interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkPosted(ctx context.Context, articleID int64) error
}

type Summarizer interface {
	Summarizer(text string) (string, error)
}

type Notifier struct {
	articlesRepository ArticlesRepository
	summarizer         Summarizer
	bot                *tgbotapi.BotAPI
	sendInterval       time.Duration
	lookupTimeWindow   time.Duration
	channelID          int64
}

func New(articlesRepository ArticlesRepository, summarizer Summarizer, bot *tgbotapi.BotAPI, sendInterval, lookupTimeWindow time.Duration, channelID int64) *Notifier {
	return &Notifier{
		articlesRepository: articlesRepository,
		summarizer:         summarizer,
		bot:                bot,
		sendInterval:       sendInterval,
		lookupTimeWindow:   lookupTimeWindow,
		channelID:          channelID,
	}
}

func (n *Notifier) Start(ctx context.Context) {
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()
	log.Println("Notifier started, will send articles every", n.sendInterval)

	if err := n.SelectAndSendArticle(ctx); err != nil {
		log.Printf("initial article selection failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				log.Printf("failed to select and send article: %v", err)
			}
		case <-ctx.Done():
			log.Println("Notifier stopped:", ctx.Err())
			return
		}
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneArticles, err := n.articlesRepository.AllNotPosted(ctx, time.Now().Add(-n.lookupTimeWindow), 1)
	if err != nil {
		return fmt.Errorf("failed to fetch articles: %w", err)
	}

	if len(topOneArticles) == 0 {
		log.Println("No articles to post at the moment")
		return nil
	}

	log.Println("Selected article for posting:", topOneArticles[0].Title)
	article := topOneArticles[0]
	summary, err := n.ExtractSummary(ctx, article)
	if err != nil {
		return fmt.Errorf("failed to extract summary: %w", err)
	}

	if err := n.sendArticle(article, summary); err != nil {
		return fmt.Errorf("failed to send article: %w", err)
	}

	log.Println("Article posted successfully:", article.Title)
	return n.articlesRepository.MarkPosted(ctx, article.ID)
}

func (n *Notifier) ExtractSummary(ctx context.Context, article model.Article) (string, error) {
	var r io.Reader

	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		response, err := http.Get(article.Link)
		if err != nil {
			log.Printf("failed to fetch article content from %s: %v", article.Link, err)
			return "", fmt.Errorf("article has no summary, tried to text from http get request but failed to fetch article content: %w", err)
		}

		defer response.Body.Close()

		r = response.Body
	}

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse article content: %w", err)
	}

	summary, err := n.summarizer.Summarizer(clearText(doc.TextContent))
	if err != nil {
		return "", fmt.Errorf("failed to summarize article content: %w", err)
	}
	return "\n\n" + summary, nil
}

var redundantNewLines = regexp.MustCompile(`\n{3,}`)

func clearText(text string) string {
	return redundantNewLines.ReplaceAllString(text, "\n")
}

func (n *Notifier) sendArticle(article model.Article, summary string) error {
	const msgFormat = "*%s*%s\n\n%s"

	// Telegram requires escaping markdown characters
	msg := tgbotapi.NewMessage(n.channelID, fmt.Sprintf(
		msgFormat,
		markup.EscapeForMarkdown(article.Title),
		markup.EscapeForMarkdown(summary),
		markup.EscapeForMarkdown(article.Link),
	))

	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
