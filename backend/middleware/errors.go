package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kubektl/v0-blog-backend/apperrors"
)

func RespondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	case errors.Is(err, apperrors.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	case errors.Is(err, apperrors.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	case errors.Is(err, apperrors.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "Conflict"})
	case errors.Is(err, apperrors.ErrInvalidTransition):
		c.JSON(http.StatusConflict, gin.H{"error": "Invalid state transition"})
	case errors.Is(err, apperrors.ErrMemberOnly):
		c.JSON(http.StatusForbidden, gin.H{"error": "This content is for members only"})
	case errors.Is(err, apperrors.ErrBadRequest):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
	c.Abort()
}

func RespondError(c *gin.Context, err error, message string) {
	if message != "" {
		switch {
		case errors.Is(err, apperrors.ErrBadRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
		default:
			RespondServiceError(c, err)
			return
		}
	} else {
		RespondServiceError(c, err)
		return
	}
	c.Abort()
}
