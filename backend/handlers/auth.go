package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kubektl/v0-blog-backend/database"
	"github.com/kubektl/v0-blog-backend/middleware"
	"github.com/kubektl/v0-blog-backend/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	var user models.User
	err = database.DB.QueryRow(context.Background(),
		`INSERT INTO users (email, password, name, avatar) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, email, name, bio, avatar, followers, created_at, updated_at`,
		req.Email, string(hashedPassword), req.Name,
		"https://api.dicebear.com/7.x/avataaars/svg?seed="+req.Name,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Bio, &user.Avatar, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT
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

	// Find user
	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow(context.Background(),
		`SELECT id, email, password, name, bio, avatar, followers, created_at, updated_at 
		 FROM users WHERE email = $1`, req.Email,
	).Scan(&user.ID, &user.Email, &hashedPassword, &user.Name, &user.Bio, &user.Avatar, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT
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
		`SELECT id, email, name, bio, avatar, followers, created_at, updated_at 
		 FROM users WHERE id = $1`, userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Bio, &user.Avatar, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

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
		`SELECT id, email, name, bio, avatar, followers, created_at, updated_at 
		 FROM users WHERE id = $1`, userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Bio, &user.Avatar, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req struct {
		Name   string `json:"name"`
		Bio    string `json:"bio"`
		Avatar string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := database.DB.QueryRow(context.Background(),
		`UPDATE users SET name = COALESCE(NULLIF($1, ''), name), 
		 bio = COALESCE($2, bio), avatar = COALESCE(NULLIF($3, ''), avatar),
		 updated_at = CURRENT_TIMESTAMP
		 WHERE id = $4
		 RETURNING id, email, name, bio, avatar, followers, created_at, updated_at`,
		req.Name, req.Bio, req.Avatar, userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Bio, &user.Avatar, &user.Followers, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
