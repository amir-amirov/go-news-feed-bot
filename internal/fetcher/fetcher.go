package fetcher

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/amir-amirov/go-news-feed-bot/internal/model"
	"github.com/amir-amirov/go-news-feed-bot/internal/source"
	"go.tomakado.io/containers/set"
)

type articlesRepository interface {
	Store(ctx context.Context, article model.Article) error
}

// Although source storage might have methods like these,
// they are not used by Fetcher, that's why interface has only Sources method.
type sourcesRepository interface {
	Sources(ctx context.Context) ([]model.Source, error)
	// SourceByID(ctx context.Context, id int64) (*model.Source, error)
	// Add(ctx context.Context, source model.Source) (int64, error)
	// Delete(ctx context.Context, id int64) error
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articlesRepository articlesRepository
	sourcesRepository  sourcesRepository

	fetchInterval  time.Duration
	filterKeywords []string
}

func New(articleRepo articlesRepository, sourceRepo sourcesRepository, fetchInterval time.Duration, filterKeywords []string) *Fetcher {
	return &Fetcher{
		articlesRepository: articleRepo,
		sourcesRepository:  sourceRepo,
		fetchInterval:      fetchInterval,
		filterKeywords:     filterKeywords,
	}
}

func (f *Fetcher) Start(ctx context.Context) {
	ticker := time.NewTicker(f.fetchInterval) // setInterval
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		log.Printf("initial fetch failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Fetcher stopped:", ctx.Err())
			return
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				log.Printf("fetch failed: %v", err)
			}
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sourcesRepository.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)

		rssSource := source.NewRSSSourceFromModel(src)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				log.Printf("[ERROR] Failed to fetch items from source %s: %v", source.Name(), err)
				return
			}

			if err := f.processItems(ctx, source, items); err != nil {
				log.Printf("[ERROR] Failed to process items from source %s: %v", source.Name(), err)
				return
			}
		}(rssSource)

	}
	wg.Wait()
	return nil
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.itemShouldBeSkipped(item) {
			continue
		}

		if err := f.articlesRepository.Store(ctx, model.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}

	}

	return nil
}

func (f *Fetcher) itemShouldBeSkipped(item model.Item) bool {
	categoriesSet := set.New(item.Categories...)

	for _, keyword := range f.filterKeywords {
		titleContainsKeyword := strings.Contains(strings.ToLower(item.Title), strings.ToLower(keyword))
		if categoriesSet.Contains(keyword) || titleContainsKeyword {
			return true
		}
	}
	return false
}
