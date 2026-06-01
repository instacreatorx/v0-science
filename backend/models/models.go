package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Locale    string    `json:"locale"`
	Followers int       `json:"followers"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Article struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Excerpt       string    `json:"excerpt"`
	Content       string    `json:"content"`
	Image         string    `json:"image,omitempty"`
	AuthorID      int64     `json:"author_id"`
	Author        *User     `json:"author,omitempty"`
	Tags          []string  `json:"tags"`
	ReadTime      string    `json:"read_time"`
	IsMemberOnly  bool      `json:"is_member_only"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Comment struct {
	ID        int64     `json:"id"`
	ArticleID int64     `json:"article_id"`
	UserID    int64     `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Like struct {
	ID        int64     `json:"id"`
	ArticleID int64     `json:"article_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Bookmark struct {
	ID        int64     `json:"id"`
	ArticleID int64     `json:"article_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Request/Response DTOs
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SendOtpRequest struct {
	Phone string `json:"phone" binding:"required"`
	Name  string `json:"name"`
}

type VerifyOtpRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required,len=6"`
	Name  string `json:"name"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	IsNewUser bool   `json:"is_new_user,omitempty"`
}

type SendOtpResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
}

type CreateArticleRequest struct {
	Title        string   `json:"title" binding:"required"`
	Excerpt      string   `json:"excerpt" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	Image        string   `json:"image"`
	Tags         []string `json:"tags"`
	IsMemberOnly bool     `json:"is_member_only"`
}

type UpdateArticleRequest struct {
	Title        string   `json:"title"`
	Excerpt      string   `json:"excerpt"`
	Content      string   `json:"content"`
	Image        string   `json:"image"`
	Tags         []string `json:"tags"`
	IsMemberOnly bool     `json:"is_member_only"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type UpdateUserRequest struct {
	Name   string `json:"name"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
	Locale string `json:"locale"`
}
