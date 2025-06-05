# Go News Feed Bot

The project is written in **Go**. It consists of two workers: **Fetcher** and **Notifier**. The first one fetches news articles from **RSS feeds** and stores them in a **PostgreSQL** database. The second one delivers articles to users posting them in the Telegram channel ([@golangnewslatest](https://t.me/golangnewslatest)). This is done by connecting to Telegram bot which is an admin in the channel. Apart from that, before posting the article is summarized using **OpenAI API**. The project follows SOLID principles and uses dependency injection for maintainability and testability. It is containerized using Docker and Docker Compose for easy deployment. An instance of **AWS EC2** was created to host the project.

## Tech Stack

- Go
- PostgreSQL
- Docker and Docker Compose (for containerized deployment)
- Telegram Bot API (for the bot admin of the Telegram channel)
- OpenAI API (for generating article summaries with `gpt-3.5-turbo`)
- RSS feeds
- AWS EC2 instance for hosting

## Architecture

The core of the Go News Feed Bot consists of two workers—**Fetcher** and **Notifier**—along with an **Admin Service**, as illustrated in the architecture diagram below (Fig. 1):

![Architecture Diagram](https://firebasestorage.googleapis.com/v0/b/auth-2c46a.appspot.com/o/news_feed_bot_diagram.png?alt=media&token=6e527ec1-4ab3-4161-b1fb-3f6bd2c6dc8a)

### Fetcher

The Fetcher worker periodically retrieves news articles from predefined RSS feeds (stored in the `sources` table). To optimize performance, it uses goroutines to fetch articles from each source concurrently. Fetched articles are then stored in the `articles` table in a PostgreSQL database (see Database Schema section).

### Notifier

The Notifier worker operates on a set interval, querying the `articles` table for entries where the `posted_at` field is `NULL` (indicating the article has not yet been posted). For each unposted article:

1. The Notifier generates a summary using the OpenAI API (model: `gpt-3.5-turbo`).
2. It constructs a post containing the article's title, summary, and link.
3. The post is sent to a Telegram Bot, which is an admin of the [@golangnewslatest](https://t.me/golangnewslatest) channel, for publication.

An example of the Telegram channel posts is shown below (Fig. 2):

![Telegram Channel Example](https://firebasestorage.googleapis.com/v0/b/auth-2c46a.appspot.com/o/Screenshot%202025-06-05%20at%2013.29.06.png?alt=media&token=a380026c-0fca-482f-a387-1f84ef2262df)

### Admin Service

The Admin Service provides a set of commands accessible via the Telegram Bot to manage news sources. These commands allow administrators to create, list, and manage sources stored in the `sources` table of the database.

## Database Schema

The bot uses a PostgreSQL database with two main tables: `sources` and `articles`. The schema definitions are as follows:

### `sources` Table

Stores information about RSS feed sources.

```sql
CREATE TABLE sources
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    feed_url   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);
```

### `articles` Table

Stores fetched news articles, linked to their respective sources.

```sql
CREATE TABLE articles
(
    id           BIGSERIAL PRIMARY KEY,
    source_id    BIGINT       NOT NULL,
    title        VARCHAR(255) NOT NULL,
    summary      TEXT         NOT NULL,
    link         TEXT         NOT NULL UNIQUE,
    published_at TIMESTAMP    NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
    posted_at    TIMESTAMP,
    CONSTRAINT fk_articles_source_id
        FOREIGN KEY (source_id)
            REFERENCES sources (id)
            ON DELETE CASCADE
);
```

## Design Principles

- **SOLID Principles**: The codebase adheres to SOLID principles to ensure maintainability, scalability, and flexibility.
- **Dependency Injection**: Dependencies are injected to promote loose coupling and improve testability.
- **Modular Design**: The bot is structured into distinct modules (Fetcher, Notifier, Admin Service) for clarity and reusability.
- **Concurrency**: Goroutines are utilized in the Fetcher to handle multiple RSS feeds efficiently.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/amir-amirov/go-news-feed-bot.git
   cd go-news-feed-bot
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   go build -o news-feed-bot
   ```
