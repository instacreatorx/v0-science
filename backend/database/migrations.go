package database

import "fmt"

// runSchemaMigrations applies idempotent SQL migrations safe for existing Railway/legacy DBs.
// Pre-existing tables are NOT passed through GORM AutoMigrate — GORM would try to drop
// constraint names like uni_users_email that never existed (legacy uses users_email_key).
func runSchemaMigrations() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE,
			password VARCHAR(255),
			name VARCHAR(255) NOT NULL,
			bio TEXT DEFAULT '',
			avatar VARCHAR(500) DEFAULT '',
			followers INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS articles (
			id SERIAL PRIMARY KEY,
			title VARCHAR(500) NOT NULL,
			excerpt TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			image VARCHAR(500) DEFAULT '',
			author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			tags TEXT[] DEFAULT '{}',
			read_time VARCHAR(50) DEFAULT '5 min read',
			is_member_only BOOLEAN DEFAULT FALSE,
			likes_count INT DEFAULT 0,
			comments_count INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS likes (
			id SERIAL PRIMARY KEY,
			article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(article_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS bookmarks (
			id SERIAL PRIMARY KEY,
			article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(article_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS follows (
			id SERIAL PRIMARY KEY,
			follower_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			following_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(follower_id, following_id)
		)`,
		`CREATE TABLE IF NOT EXISTS otp_codes (
			id SERIAL PRIMARY KEY,
			phone VARCHAR(20) NOT NULL,
			code VARCHAR(6) NOT NULL,
			name VARCHAR(255) DEFAULT '',
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS locale VARCHAR(10) DEFAULT 'en'`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(32) DEFAULT 'user'`,
		`ALTER TABLE users ALTER COLUMN email DROP NOT NULL`,
		`ALTER TABLE users ALTER COLUMN password DROP NOT NULL`,

		`ALTER TABLE articles ADD COLUMN IF NOT EXISTS slug VARCHAR(500) DEFAULT ''`,
		`ALTER TABLE articles ADD COLUMN IF NOT EXISTS status VARCHAR(32) DEFAULT 'draft'`,
		`ALTER TABLE articles ADD COLUMN IF NOT EXISTS published_at TIMESTAMP`,
		`ALTER TABLE articles ADD COLUMN IF NOT EXISTS team_id INT`,

		`CREATE INDEX IF NOT EXISTS idx_articles_author_id ON articles(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_slug ON articles(slug)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_article_id ON comments(article_id)`,
		`CREATE INDEX IF NOT EXISTS idx_likes_article_id ON likes(article_id)`,
		`CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_otp_codes_phone ON otp_codes(phone)`,
	}

	for _, stmt := range statements {
		if err := DB.Exec(stmt).Error; err != nil {
			return fmt.Errorf("migration failed: %w\nSQL: %s", err, stmt)
		}
	}

	// New tables (never use GORM AutoMigrate — it re-migrates users via Team FK)
	newTables := []string{
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash VARCHAR(64) NOT NULL UNIQUE,
			expires_at TIMESTAMP NOT NULL,
			revoked_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
		`CREATE TABLE IF NOT EXISTS teams (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			bio TEXT DEFAULT '',
			avatar VARCHAR(500) DEFAULT '',
			owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			verified_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_teams_owner_id ON teams(owner_id)`,
		`CREATE TABLE IF NOT EXISTS team_members (
			id SERIAL PRIMARY KEY,
			team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(32) NOT NULL DEFAULT 'writer',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(team_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS team_verification_requests (
			id SERIAL PRIMARY KEY,
			team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
			submitted_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			proof_text TEXT NOT NULL,
			proof_url VARCHAR(500) DEFAULT '',
			status VARCHAR(32) DEFAULT 'pending',
			reviewed_by INT REFERENCES users(id) ON DELETE SET NULL,
			rejection_reason TEXT DEFAULT '',
			reviewed_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_team_verification_team_id ON team_verification_requests(team_id)`,
		`CREATE INDEX IF NOT EXISTS idx_team_verification_status ON team_verification_requests(status)`,
	}
	for _, stmt := range newTables {
		if err := DB.Exec(stmt).Error; err != nil {
			return fmt.Errorf("new table migration failed: %w\nSQL: %s", err, stmt)
		}
	}

	// Phone unique — partial index (nullable phone); skip if legacy users_phone_key exists
	if err := DB.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'users_phone_key'
			) AND NOT EXISTS (
				SELECT 1 FROM pg_indexes WHERE indexname = 'users_phone_key'
			) THEN
				CREATE UNIQUE INDEX users_phone_key ON users (phone) WHERE phone IS NOT NULL;
			END IF;
		END $$;
	`).Error; err != nil {
		return fmt.Errorf("migration failed (users_phone_key): %w", err)
	}

	return nil
}

func ensureSlugUniqueIndex() error {
	// After backfill assigns slugs, enforce uniqueness
	return DB.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_indexes WHERE indexname = 'articles_slug_key'
			) THEN
				CREATE UNIQUE INDEX articles_slug_key ON articles (slug) WHERE slug <> '';
			END IF;
		END $$;
	`).Error
}
