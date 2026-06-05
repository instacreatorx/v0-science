package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/config"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var phoneRegex = regexp.MustCompile(`^\+?[\d]{10,15}$`)

type AuthService struct {
	cfg      *config.Config
	users    repositories.UserRepository
	auth     repositories.AuthRepository
	jwtSecret []byte
}

func NewAuthService(cfg *config.Config, users repositories.UserRepository, auth repositories.AuthRepository) *AuthService {
	return &AuthService{
		cfg:       cfg,
		users:     users,
		auth:      auth,
		jwtSecret: []byte(cfg.JWTSecret),
	}
}

func (s *AuthService) JWTSecret() []byte {
	return s.jwtSecret
}

func normalizePhone(phone string) string {
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	return cleaned
}

func (s *AuthService) generateOTP() (string, error) {
	if s.cfg.DevFixedOTP {
		return "123456", nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func (s *AuthService) SendOTP(phone, name string) (*models.SendOtpResponse, error) {
	phone = normalizePhone(phone)
	if !phoneRegex.MatchString(phone) {
		return nil, apperrors.ErrBadRequest
	}

	code, err := s.generateOTP()
	if err != nil {
		return nil, err
	}

	_ = s.auth.DeleteOTPsByPhone(phone)
	if err := s.auth.CreateOTP(&models.OtpCode{
		Phone:     phone,
		Code:      code,
		Name:      name,
		ExpiresAt: time.Now().Add(s.cfg.OTPTTL),
	}); err != nil {
		return nil, err
	}

	log.Printf("[OTP] Phone: %s, Code: %s", phone, code)
	return &models.SendOtpResponse{Message: "Verification code sent", ExpiresIn: 60}, nil
}

func (s *AuthService) VerifyOTP(phone, code, name string) (*models.AuthResponse, error) {
	phone = normalizePhone(phone)
	otp, err := s.auth.FindLatestOTP(phone)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	if time.Now().After(otp.ExpiresAt) {
		return nil, apperrors.ErrUnauthorized
	}
	if code != otp.Code {
		return nil, apperrors.ErrUnauthorized
	}
	_ = s.auth.DeleteOTPByPhone(phone)

	user, isNew, err := s.findOrCreatePhoneUser(phone, name, otp.Name)
	if err != nil {
		return nil, err
	}
	return s.issueTokens(user, isNew)
}

func (s *AuthService) findOrCreatePhoneUser(phone, reqName, otpName string) (*models.User, bool, error) {
	user, err := s.users.FindByPhone(phone)
	if err == nil {
		return user, false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, false, err
	}

	name := reqName
	if name == "" {
		name = otpName
	}
	if name == "" {
		name = "User"
	}

	phoneCopy := phone
	newUser := &models.User{
		Phone:  &phoneCopy,
		Name:   name,
		Avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=" + name,
		Locale: "en",
		Role:   models.RoleUser,
	}
	if err := s.users.Create(newUser); err != nil {
		return nil, false, err
	}
	return newUser, true, nil
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	_, err := s.users.FindByEmail(req.Email)
	if err == nil {
		return nil, apperrors.ErrConflict
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	email := req.Email
	user := &models.User{
		Email:    &email,
		Password: string(hashed),
		Name:     req.Name,
		Avatar:   "https://api.dicebear.com/7.x/avataaars/svg?seed=" + req.Name,
		Locale:   "en",
		Role:     models.RoleUser,
	}
	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	return s.issueTokens(user, false)
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.users.FindByEmail(req.Email)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	return s.issueTokens(user, false)
}

func (s *AuthService) Refresh(refreshToken string) (*models.AuthResponse, error) {
	hash := hashToken(refreshToken)
	stored, err := s.auth.FindRefreshTokenByHash(hash)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	if stored.RevokedAt != nil || time.Now().After(stored.ExpiresAt) {
		return nil, apperrors.ErrUnauthorized
	}

	user, err := s.users.FindByID(stored.UserID)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	_ = s.auth.RevokeRefreshToken(stored.ID)
	return s.issueTokens(user, false)
}

func (s *AuthService) Logout(refreshToken string) error {
	hash := hashToken(refreshToken)
	stored, err := s.auth.FindRefreshTokenByHash(hash)
	if err != nil {
		return apperrors.ErrUnauthorized
	}
	return s.auth.RevokeRefreshToken(stored.ID)
}

func (s *AuthService) issueTokens(user *models.User, isNew bool) (*models.AuthResponse, error) {
	access, err := s.createAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshRaw, refreshHash, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.auth.CreateRefreshToken(&models.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(s.cfg.RefreshTTL),
	}); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:        access,
		RefreshToken: refreshRaw,
		User:         *user,
		IsNewUser:    isNew,
	}, nil
}

func (s *AuthService) createAccessToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.cfg.AccessTTL).Unix(),
	})
	return token.SignedString(s.jwtSecret)
}

func generateRefreshToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	raw = hex.EncodeToString(b)
	return raw, hashToken(raw), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func ParseUserIDFromToken(tokenString string, secret []byte) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return 0, apperrors.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, apperrors.ErrUnauthorized
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, apperrors.ErrUnauthorized
	}
	return int64(userID), nil
}
