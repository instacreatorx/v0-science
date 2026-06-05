package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/services"
)

type AuthMiddleware struct {
	auth *services.AuthService
}

func NewAuthMiddleware(auth *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := m.parseToken(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID, ok := m.parseToken(c); ok {
			c.Set("user_id", userID)
		}
		c.Next()
	}
}

func (m *AuthMiddleware) SuperAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := m.parseToken(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}

func (m *AuthMiddleware) parseToken(c *gin.Context) (int64, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, false
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return 0, false
	}
	userID, err := services.ParseUserIDFromToken(tokenString, m.auth.JWTSecret())
	if err != nil {
		return 0, false
	}
	return userID, true
}

func GetOptionalUserID(c *gin.Context) int64 {
	v, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	id, ok := v.(int64)
	if !ok {
		return 0
	}
	return id
}

func RequireSuperAdmin(userService func(int64) (*models.User, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := userService(c.GetInt64("user_id"))
		if err != nil || user.Role != models.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Super admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
