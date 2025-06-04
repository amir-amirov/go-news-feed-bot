package source_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/amir-amirov/go-news-feed-bot/internal/source"
)

func TestFetch_RealRSSFeed(t *testing.T) {

	// src := source.NewRSSSourceFromModel(model.Source{
	// 	ID:      1,
	// 	Name:    "BBC News",
	// 	FeedURL: "https://feeds.bbci.co.uk/news/rss.xml",
	// })

	// src := source.RSSSource{
	// 	URL:        "https://reactnative.dev/blog/rss.xml",
	// 	SourceID:   1,
	// 	SourceName: "BBC News",
	// }
	src := source.RSSSource{
		URL:        "https://thisweekinreact.com/rss.xml",
		SourceID:   1,
		SourceName: "BBC News",
	}

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

/*
Results:

amir@Amirs-MacBook-Air source % go test
Fetched 38 items from BBC News:
Title: Defence review to send 'message to Moscow', says Healey
Link: https://www.bbc.com/news/articles/cq69vqpp2l4o
Summary: At least six new factories will be built and up to 7,000 UK-built long-range weapons will be procured.
Date: Sun, 01 Jun 2025 05:00:50 GMT
Categories: []

Title: Two dead and hundreds arrested in France after PSG Champions League win
Link: https://www.bbc.com/news/articles/ckgqyg325gno
Summary: A 17-year-old boy was stabbed in Dax and a 23-year-old man hit by a vehicle of supporters while riding a scooter in Paris.
Date: Sun, 01 Jun 2025 09:47:46 GMT
Categories: []

Title: Briton accused of plot to export US military tech
Link: https://www.bbc.com/news/articles/c0qg4q87p1zo
Summary: John Miller and Chinese man Cui Guanghai are in custody in Serbia facing a US extradition request.
Date: Sun, 01 Jun 2025 04:52:58 GMT
Categories: []

Title: Ranganathan opens up about mental health struggle
Link: https://www.bbc.com/news/articles/cy8np7zzdl3o
Summary: The comedian says he hopes to destigmatise mental health issues by being open about his own experiences.
Date: Sun, 01 Jun 2025 00:53:41 GMT
Categories: []

Title: Disposable vape ban begins - but will it have an impact?
Link: https://www.bbc.com/news/articles/c80kxx2xr77o
Summary: Retailers will no longer be able to sell single-use vapes as new laws come into force across the UK.
Date: Sat, 31 May 2025 23:12:11 GMT
Categories: []

PASS
ok      github.com/amir-amirov/go-news-feed-bot/internal/source 1.393s
*/
