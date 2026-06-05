package models

import "time"

type Comment struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ArticleID int64     `gorm:"index" json:"article_id"`
	UserID    int64     `json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Like struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ArticleID int64     `gorm:"uniqueIndex:idx_like_article_user" json:"article_id"`
	UserID    int64     `gorm:"uniqueIndex:idx_like_article_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Bookmark struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ArticleID int64     `gorm:"uniqueIndex:idx_bookmark_article_user" json:"article_id"`
	UserID    int64     `gorm:"uniqueIndex:idx_bookmark_article_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Follow struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	FollowerID  int64     `gorm:"uniqueIndex:idx_follow_pair" json:"follower_id"`
	FollowingID int64     `gorm:"uniqueIndex:idx_follow_pair" json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}
