package models

import "time"

type OtpCode struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"index" json:"phone"`
	Code      string    `json:"-"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	UserID    int64      `gorm:"index" json:"user_id"`
	TokenHash string     `gorm:"uniqueIndex:refresh_tokens_token_hash_key" json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
