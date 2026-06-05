package models

import (
	"time"

	"github.com/lib/pq"
)

const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
	ArticleStatusArchived  = "archived"
)

type Article struct {
	ID            int64          `gorm:"primaryKey" json:"id"`
	Title         string         `json:"title"`
	Slug          string         `gorm:"uniqueIndex:articles_slug_key" json:"slug"`
	Excerpt       string         `json:"excerpt"`
	Content       string         `json:"content"`
	Image         string         `json:"image,omitempty"`
	AuthorID      int64          `gorm:"index" json:"author_id"`
	Author        *User          `gorm:"foreignKey:AuthorID;constraint:-" json:"author,omitempty"`
	TeamID        *int64         `gorm:"index" json:"team_id,omitempty"`
	Team          *Team          `gorm:"foreignKey:TeamID;constraint:-" json:"team,omitempty"`
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags"`
	ReadTime      string         `gorm:"default:'5 min read'" json:"read_time"`
	Status        string         `gorm:"default:draft;index" json:"status"`
	IsMemberOnly  bool           `gorm:"default:false" json:"is_member_only"`
	LikesCount    int            `gorm:"default:0" json:"likes_count"`
	CommentsCount int            `gorm:"default:0" json:"comments_count"`
	PublishedAt   *time.Time     `json:"published_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`

	LikedByMe      *bool `gorm:"-" json:"liked_by_me,omitempty"`
	BookmarkedByMe *bool `gorm:"-" json:"bookmarked_by_me,omitempty"`
}

func (Article) TableName() string { return "articles" }
