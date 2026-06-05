package repositories

import (
	"time"

	"github.com/kubektl/v0-blog-backend/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	DeleteOTPsByPhone(phone string) error
	CreateOTP(otp *models.OtpCode) error
	FindLatestOTP(phone string) (*models.OtpCode, error)
	DeleteOTPByPhone(phone string) error
	CreateRefreshToken(token *models.RefreshToken) error
	FindRefreshTokenByHash(hash string) (*models.RefreshToken, error)
	RevokeRefreshToken(id int64) error
	RevokeAllUserTokens(userID int64) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) DeleteOTPsByPhone(phone string) error {
	return r.db.Where("phone = ?", phone).Delete(&models.OtpCode{}).Error
}

func (r *authRepository) CreateOTP(otp *models.OtpCode) error {
	return r.db.Create(otp).Error
}

func (r *authRepository) FindLatestOTP(phone string) (*models.OtpCode, error) {
	var otp models.OtpCode
	err := r.db.Where("phone = ?", phone).Order("created_at DESC").First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *authRepository) DeleteOTPByPhone(phone string) error {
	return r.db.Where("phone = ?", phone).Delete(&models.OtpCode{}).Error
}

func (r *authRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *authRepository) FindRefreshTokenByHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Where("token_hash = ?", hash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *authRepository) RevokeRefreshToken(id int64) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked_at", now).Error
}

func (r *authRepository) RevokeAllUserTokens(userID int64) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}
