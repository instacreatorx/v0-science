package database

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kubektl/v0-blog-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dsn string, ginMode string) error {
	logLevel := logger.Warn
	if ginMode != "release" {
		logLevel = logger.Info
	}

	gormLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             2 * time.Second,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	DB = db
	return nil
}

func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func Migrate() error {
	// All schema changes via idempotent SQL — never GORM AutoMigrate on Railway.
	// AutoMigrate walks related models (Team → User) and breaks legacy unique constraints.
	if err := runSchemaMigrations(); err != nil {
		return err
	}

	// Backfill slugs/status on existing rows, then add slug unique index
	if err := backfillLegacyData(); err != nil {
		return err
	}

	if err := ensureSlugUniqueIndex(); err != nil {
		return fmt.Errorf("slug index migration failed: %w", err)
	}

	return nil
}

func backfillLegacyData() error {
	var articles []models.Article
	if err := DB.Where("slug = '' OR slug IS NULL OR status = '' OR status IS NULL").Find(&articles).Error; err != nil {
		return err
	}

	for _, a := range articles {
		updates := map[string]interface{}{}
		if a.Status == "" {
			updates["status"] = models.ArticleStatusPublished
		}
		if a.Slug == "" {
			base := slugify(a.Title)
			if base == "" {
				base = fmt.Sprintf("article-%d", a.ID)
			}
			slug := base
			for i := 0; ; i++ {
				var count int64
				DB.Model(&models.Article{}).Where("slug = ? AND id != ?", slug, a.ID).Count(&count)
				if count == 0 {
					break
				}
				slug = fmt.Sprintf("%s-%d", base, i+1)
			}
			updates["slug"] = slug
		}
		if (a.Status == models.ArticleStatusPublished || updates["status"] == models.ArticleStatusPublished) && a.PublishedAt == nil {
			t := a.CreatedAt
			updates["published_at"] = t
		}
		if len(updates) > 0 {
			if err := DB.Model(&models.Article{}).Where("id = ?", a.ID).Updates(updates).Error; err != nil {
				return err
			}
		}
	}

	return DB.Model(&models.User{}).Where("role = '' OR role IS NULL").Update("role", models.RoleUser).Error
}

var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.Trim(slugRegex.ReplaceAllString(s, "-"), "-")
}

func Slugify(s string) string {
	return slugify(s)
}

func UniqueSlug(base string, excludeID int64) (string, error) {
	base = slugify(base)
	if base == "" {
		base = "untitled"
	}
	slug := base
	for i := 0; ; i++ {
		var count int64
		q := DB.Model(&models.Article{}).Where("slug = ?", slug)
		if excludeID > 0 {
			q = q.Where("id != ?", excludeID)
		}
		if err := q.Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, i+1)
	}
}

func UniqueTeamSlug(base string, excludeID int64) (string, error) {
	base = slugify(base)
	if base == "" {
		base = "team"
	}
	slug := base
	for i := 0; ; i++ {
		var count int64
		q := DB.Model(&models.Team{}).Where("slug = ?", slug)
		if excludeID > 0 {
			q = q.Where("id != ?", excludeID)
		}
		if err := q.Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, i+1)
	}
}
