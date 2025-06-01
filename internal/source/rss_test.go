package source_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/amir-amirov/go-news-feed-bot/internal/model"
	"github.com/amir-amirov/go-news-feed-bot/internal/source"
)

func TestFetch_RealRSSFeed(t *testing.T) {

	src := source.NewRSSSourceFromModel(model.Source{
		ID:      1,
		Name:    "BBC News",
		FeedURL: "https://feeds.bbci.co.uk/news/rss.xml",
	})

	// src := source.RSSSource{
	// 	URL:        "https://feeds.bbci.co.uk/news/rss.xml", // BBC News RSS feed
	// 	SourceID:   1,
	// 	SourceName: "BBC News",
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := src.Fetch(ctx)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	fmt.Printf("Fetched %d items from %s:\n", len(items), src.SourceName)
	for i, item := range items {
		if i >= 5 {
			break // print only first 5 items
		}
		fmt.Printf("Title: %s\nLink: %s\nSummary: %s\nDate: %s\nCategories: %v\n\n",
			item.Title, item.Link, item.Summary, item.Date.Format(time.RFC1123), item.Categories)
	}
}
