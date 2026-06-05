package models

import "time"

const (
	RoleUser       = "user"
	RoleSuperAdmin = "super_admin"
)

type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Email     *string   `gorm:"uniqueIndex:users_email_key" json:"email,omitempty"`
	Phone     *string   `gorm:"uniqueIndex:users_phone_key" json:"phone,omitempty"`
	Password  string    `gorm:"column:password" json:"-"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Locale    string    `gorm:"default:en" json:"locale"`
	Role      string    `gorm:"default:user" json:"role"`
	Followers int       `gorm:"default:0" json:"followers"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsFollowing *bool `gorm:"-" json:"is_following,omitempty"`
}

func (User) TableName() string { return "users" }
