package models

import "time"

const (
	TeamRoleOwner  = "owner"
	TeamRoleAdmin  = "admin"
	TeamRoleEditor = "editor"
	TeamRoleWriter = "writer"

	VerificationStatusPending  = "pending"
	VerificationStatusApproved = "approved"
	VerificationStatusRejected = "rejected"
)

type Team struct {
	ID         int64      `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	Slug       string     `gorm:"uniqueIndex:teams_slug_key" json:"slug"`
	Bio        string     `json:"bio,omitempty"`
	Avatar     string     `json:"avatar,omitempty"`
	OwnerID    int64      `gorm:"index" json:"owner_id"`
	Owner      *User      `gorm:"foreignKey:OwnerID;constraint:-" json:"owner,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	IsVerified bool       `gorm:"-" json:"is_verified"`
}

type TeamMember struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TeamID    int64     `gorm:"uniqueIndex:idx_team_member" json:"team_id"`
	UserID    int64     `gorm:"uniqueIndex:idx_team_member" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID;constraint:-" json:"user,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type TeamVerificationRequest struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	TeamID          int64      `gorm:"index" json:"team_id"`
	Team            *Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	SubmittedBy     int64      `json:"submitted_by"`
	ProofText       string     `json:"proof_text"`
	ProofURL        string     `json:"proof_url,omitempty"`
	Status          string     `gorm:"default:pending;index" json:"status"`
	ReviewedBy      *int64     `json:"reviewed_by,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
