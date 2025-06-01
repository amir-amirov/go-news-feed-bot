package db

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
