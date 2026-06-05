package models

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

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
	IsNewUser    bool   `json:"is_new_user,omitempty"`
}

type SendOtpResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
}

type CreateArticleRequest struct {
	Title        string   `json:"title" binding:"required"`
	Excerpt      string   `json:"excerpt"`
	Content      string   `json:"content"`
	Image        string   `json:"image"`
	Tags         []string `json:"tags"`
	TeamID       *int64   `json:"team_id"`
	IsMemberOnly bool     `json:"is_member_only"`
}

type UpdateArticleRequest struct {
	Title        string   `json:"title"`
	Excerpt      string   `json:"excerpt"`
	Content      string   `json:"content"`
	Image        string   `json:"image"`
	Tags         []string `json:"tags"`
	TeamID       *int64   `json:"team_id"`
	IsMemberOnly *bool    `json:"is_member_only"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateUserRequest struct {
	Name   string `json:"name"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
	Locale string `json:"locale"`
}

type CreateTeamRequest struct {
	Name   string `json:"name" binding:"required"`
	Slug   string `json:"slug" binding:"required"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
}

type UpdateTeamRequest struct {
	Name   string `json:"name"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
}

type AddTeamMemberRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type TeamVerifyRequest struct {
	ProofText string `json:"proof_text" binding:"required"`
	ProofURL  string `json:"proof_url"`
}

type RejectVerificationRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
