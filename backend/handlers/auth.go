package handlers

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kubektl/v0-blog-backend/database"
	"github.com/kubektl/v0-blog-backend/middleware"
	"github.com/kubektl/v0-blog-backend/models"
	"golang.org/x/crypto/bcrypt"
)

var phoneRegex = regexp.MustCompile(`^\+?[\d]{10,15}$`)

func normalizePhone(phone string) string {
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	return cleaned
}

func generateOTP() (string, error) {
	// In development, use fixed OTP for easier testing
	if os.Getenv("GIN_MODE") != "release" {
		return "123456", nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func SendOTP(c *gin.Context) {
	var req models.SendOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phone := normalizePhone(req.Phone)
	if !phoneRegex.MatchString(phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		return
	}

	code, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	expiresAt := time.Now().Add(5 * time.Minute)

	// Remove existing OTPs for this phone
	_, _ = database.DB.Exec(context.Background(),
		"DELETE FROM otp_codes WHERE phone = $1", phone)

	_, err = database.DB.Exec(context.Background(),
		`INSERT INTO otp_codes (phone, code, name, expires_at) VALUES ($1, $2, $3, $4)`,
		phone, code, req.Name, expiresAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
		return
	}

	// In production, send SMS here. For dev, log the OTP.
	log.Printf("[OTP] Phone: %s, Code: %s", phone, code)

	c.JSON(http.StatusOK, models.SendOtpResponse{
		Message:   "Verification code sent",
		ExpiresIn: 60,
	})
}

func VerifyOTP(c *gin.Context) {
	var req models.VerifyOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phone := normalizePhone(req.Phone)

	var storedCode, storedName string
	var expiresAt time.Time
	err := database.DB.QueryRow(context.Background(),
		`SELECT code, name, expires_at FROM otp_codes 
		 WHERE phone = $1 ORDER BY created_at DESC LIMIT 1`, phone,
	).Scan(&storedCode, &storedName, &expiresAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification code"})
		return
	}

	if time.Now().After(expiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Verification code expired"})
		return
	}

	if req.Code != storedCode {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
		return
	}

	// Delete used OTP
	_, _ = database.DB.Exec(context.Background(), "DELETE FROM otp_codes WHERE phone = $1", phone)

	// Check if user exists
	var user models.User
	var isNewUser bool
	err = database.DB.QueryRow(context.Background(),
		`SELECT id, COALESCE(email, ''), COALESCE(phone, ''), name, COALESCE(bio, ''), 
		 COALESCE(avatar, ''), COALESCE(locale, 'en'), followers, created_at, updated_at 
		 FROM users WHERE phone = $1`, phone,
	).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Bio, &user.Avatar,
		&user.Locale, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Create new user
		name := req.Name
		if name == "" {
			name = storedName
		}
		if name == "" {
			name = "User"
		}

		isNewUser = true
		err = database.DB.QueryRow(context.Background(),
			`INSERT INTO users (phone, name, avatar, locale) 
			 VALUES ($1, $2, $3, 'en') 
			 RETURNING id, COALESCE(email, ''), COALESCE(phone, ''), name, COALESCE(bio, ''), 
			 COALESCE(avatar, ''), COALESCE(locale, 'en'), followers, created_at, updated_at`,
			phone, name,
			"https://api.dicebear.com/7.x/avataaars/svg?seed="+name,
		).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Bio, &user.Avatar,
			&user.Locale, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * 7 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(middleware.GetJWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:     tokenString,
		User:      user,
		IsNewUser: isNewUser,
	})
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err := database.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var user models.User
	err = database.DB.QueryRow(context.Background(),
		`INSERT INTO users (email, password, name, avatar, locale) 
		 VALUES ($1, $2, $3, $4, 'en') 
		 RETURNING id, COALESCE(email, ''), COALESCE(phone, ''), name, COALESCE(bio, ''), 
		 COALESCE(avatar, ''), COALESCE(locale, 'en'), followers, created_at, updated_at`,
		req.Email, string(hashedPassword), req.Name,
		"https://api.dicebear.com/7.x/avataaars/svg?seed="+req.Name,
	).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Bio, &user.Avatar,
		&user.Locale, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * 7 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(middleware.GetJWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: tokenString,
		User:  user,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow(context.Background(),
		`SELECT id, COALESCE(email, ''), COALESCE(phone, ''), password, name, COALESCE(bio, ''), 
		 COALESCE(avatar, ''), COALESCE(locale, 'en'), followers, created_at, updated_at 
		 FROM users WHERE email = $1`, req.Email,
	).Scan(&user.ID, &user.Email, &user.Phone, &hashedPassword, &user.Name, &user.Bio,
		&user.Avatar, &user.Locale, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * 7 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(middleware.GetJWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: tokenString,
		User:  user,
	})
}

func GetCurrentUser(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var user models.User
	err := database.DB.QueryRow(context.Background(),
		`SELECT id, COALESCE(email, ''), COALESCE(phone, ''), name, COALESCE(bio, ''), 
		 COALESCE(avatar, ''), COALESCE(locale, 'en'), followers, created_at, updated_at 
		 FROM users WHERE id = $1`, userID,
	).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Bio, &user.Avatar,
		&user.Locale, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	err := database.DB.QueryRow(context.Background(),
		`SELECT id, COALESCE(email, ''), COALESCE(phone, ''), name, COALESCE(bio, ''), 
		 COALESCE(avatar, ''), COALESCE(locale, 'en'), followers, created_at, updated_at 
		 FROM users WHERE id = $1`, userID,
	).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Bio, &user.Avatar,
		&user.Locale, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate locale if provided
	if req.Locale != "" && req.Locale != "en" && req.Locale != "fa" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid locale"})
		return
	}

	var user models.User
	err := database.DB.QueryRow(context.Background(),
		`UPDATE users SET 
		 name = COALESCE(NULLIF($1, ''), name), 
		 bio = COALESCE($2, bio), 
		 avatar = COALESCE(NULLIF($3, ''), avatar),
		 locale = COALESCE(NULLIF($4, ''), locale),
		 updated_at = CURRENT_TIMESTAMP
		 WHERE id = $5
		 RETURNING id, COALESCE(email, ''), COALESCE(phone, ''), name, COALESCE(bio, ''), 
		 COALESCE(avatar, ''), COALESCE(locale, 'en'), followers, created_at, updated_at`,
		req.Name, req.Bio, req.Avatar, req.Locale, userID,
	).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Bio, &user.Avatar,
		&user.Locale, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
