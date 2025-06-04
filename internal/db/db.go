package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connStr string) {

	// connStr := "host=localhost port=5430 user=postgres password=postgress dbname=news_feed_bot sslmode=disable"
	var err error

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	err = DB.Ping()
	if err != nil {
		panic("Failed to ping the database: " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	fmt.Println("Successfully connected to PostgreSQL!")
}

/*
CREATE TABLE IF NOT EXISTS sources(
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	feed_url VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)

DROP TABLE IF EXISTS sources

// ----------------------------------------------------------

CREATE TABLE IF NOT EXISTS articles(
	id SERIAL PRIMARY KEY,
	source_id INT NOT NULL,
	title VARCHAR(255) NOT NULL,
	link VARCHAR(255) NOT NULL,
	summary TEXT NOT NULL,
	published_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	posted_at TIMESTAMP NOT NULL DEFAULT NOW(),
	CONSTRAINT fk_articles_sources_id
		FOREIGN KEY (source_id) REFERENCES sources(id)
)

DROP TABLE IF EXISTS articles
*/
