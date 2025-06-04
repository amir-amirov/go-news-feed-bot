package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/amir-amirov/go-news-feed-bot/internal/model"
)

type ArticlePostgresStorage struct {
	db *sql.DB
}

func NewArticlePostgresStorage(db *sql.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{
		db: db,
	}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	query := `
		INSERT INTO articles(source_id, title, link, summary, published_at)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (link) DO NOTHING
	`
	// The ON CONFLICT clause ensures that if an article with the same link already exists,
	// it will not be inserted and it will not return an error.

	_, err := s.db.ExecContext(ctx, query, article.SourceID, article.Title, article.Link, article.Summary, article.PublishedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	log.Println("Fetching all not posted articles since", since, "with limit", limit)
	query := `
		SELECT id, source_id, title, link, summary, published_at, posted_at, created_at
		FROM articles
		WHERE posted_at IS NULL AND published_at >= $1::timestamp
		ORDER BY published_at DESC
		LIMIT $2
	`
	rows, err := s.db.QueryContext(ctx, query, since.UTC().Format(time.RFC3339), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var dbArticle dbArticle
		if err := rows.Scan(&dbArticle.ID, &dbArticle.SourceID, &dbArticle.Title, &dbArticle.Link, &dbArticle.Summary, &dbArticle.PublishedAt, &dbArticle.PostedAt, &dbArticle.CreatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, *modelArticleFromDB(dbArticle))
	}

	return articles, nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	query := `
		UPDATE articles
		SET posted_at = $1::timestamp
		WHERE id = $2
	`

	_, err := s.db.ExecContext(ctx, query, time.Now().UTC().Format(time.RFC3339), id)

	if err != nil {
		return err
	}
	return nil
}

type dbArticle struct {
	ID          int64        `db:"id"`
	SourceID    int64        `db:"source_id"`
	Title       string       `db:"title"`
	Link        string       `db:"link"`
	Summary     string       `db:"summary"`
	PublishedAt time.Time    `db:"published_at"`
	PostedAt    sql.NullTime `db:"posted_at"`
	CreatedAt   time.Time    `db:"created_at"`
}

func modelArticleFromDB(dbArticle dbArticle) *model.Article {
	return &model.Article{
		ID:          dbArticle.ID,
		SourceID:    dbArticle.SourceID,
		Title:       dbArticle.Title,
		Link:        dbArticle.Link,
		Summary:     dbArticle.Summary,
		PublishedAt: dbArticle.PublishedAt,
		CreatedAt:   dbArticle.CreatedAt,
	}
}
