package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	if err := DB.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func Migrate() error {
	ctx := context.Background()

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		bio TEXT DEFAULT '',
		avatar VARCHAR(500) DEFAULT '',
		followers INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS articles (
		id SERIAL PRIMARY KEY,
		title VARCHAR(500) NOT NULL,
		excerpt TEXT NOT NULL,
		content TEXT NOT NULL,
		image VARCHAR(500) DEFAULT '',
		author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		tags TEXT[] DEFAULT '{}',
		read_time VARCHAR(50) DEFAULT '5 min read',
		is_member_only BOOLEAN DEFAULT FALSE,
		likes_count INT DEFAULT 0,
		comments_count INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS likes (
		id SERIAL PRIMARY KEY,
		article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(article_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS bookmarks (
		id SERIAL PRIMARY KEY,
		article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(article_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS follows (
		id SERIAL PRIMARY KEY,
		follower_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		following_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(follower_id, following_id)
	);

	CREATE INDEX IF NOT EXISTS idx_articles_author_id ON articles(author_id);
	CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_comments_article_id ON comments(article_id);
	CREATE INDEX IF NOT EXISTS idx_likes_article_id ON likes(article_id);
	CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
	`

	_, err := DB.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Additive migrations for phone auth and locale
	additiveMigrations := []string{
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20) UNIQUE`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS locale VARCHAR(10) DEFAULT 'en'`,
		`ALTER TABLE users ALTER COLUMN email DROP NOT NULL`,
		`ALTER TABLE users ALTER COLUMN password DROP NOT NULL`,
		`CREATE TABLE IF NOT EXISTS otp_codes (
			id SERIAL PRIMARY KEY,
			phone VARCHAR(20) NOT NULL,
			code VARCHAR(6) NOT NULL,
			name VARCHAR(255) DEFAULT '',
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_otp_codes_phone ON otp_codes(phone)`,
	}

	for _, migration := range additiveMigrations {
		if _, err := DB.Exec(ctx, migration); err != nil {
			return fmt.Errorf("failed to run additive migration: %w", err)
		}
	}

	return nil
}
