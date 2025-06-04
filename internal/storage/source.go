package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/amir-amirov/go-news-feed-bot/internal/model"
	_ "github.com/lib/pq"
)

type SourcePostgresStorage struct {
	db *sql.DB
}

func NewSourcePostgresStorage(db *sql.DB) *SourcePostgresStorage {
	return &SourcePostgresStorage{
		db: db,
	}
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	query := `SELECT * FROM sources`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sources []model.Source
	for rows.Next() {
		var dbSrc dbSource
		if err := rows.Scan(&dbSrc.ID, &dbSrc.Name, &dbSrc.FeedURL, &dbSrc.CreatedAt); err != nil {
			return nil, err
		}
		sources = append(sources, *modelSourceFromDB(dbSrc))
	}

	return sources, nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int64) (*model.Source, error) {
	// query := `SELECT * FROM sources WHERE id = $1` // not recommended
	// Since we are scanning into a dbSource struct, we'd use a more specific query
	query := `SELECT id, name, feed_url, created_at FROM sources WHERE id = $1`

	var source dbSource

	row := s.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&source.ID, &source.Name, &source.FeedURL, &source.CreatedAt); err != nil {
		return nil, err
	}

	return modelSourceFromDB(source), nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	query := `
		INSERT INTO sources (name, feed_url)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	err := s.db.QueryRowContext(ctx, query, source.Name, source.FeedURL).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	query := `DELETTE FROM source WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}

func modelSourceFromDB(dbSource dbSource) *model.Source {
	return &model.Source{
		ID:        dbSource.ID,
		Name:      dbSource.Name,
		FeedURL:   dbSource.FeedURL,
		CreatedAt: dbSource.CreatedAt,
	}
}
