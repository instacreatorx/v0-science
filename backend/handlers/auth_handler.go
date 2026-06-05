package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/middleware"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/services"
)

type AuthHandler struct {
	auth  *services.AuthService
	users *services.UserService
}

func NewAuthHandler(auth *services.AuthService, users *services.UserService) *AuthHandler {
	return &AuthHandler{auth: auth, users: users}
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req models.SendOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.auth.SendOTP(req.Phone, req.Name)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.auth.VerifyOTP(req.Phone, req.Code, req.Name)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.auth.Register(req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.auth.Login(req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.auth.Refresh(req.RefreshToken)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	if err := h.auth.Logout(req.RefreshToken); err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, err := h.users.GetCurrent(c.GetInt64("user_id"))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	user, err := h.users.Update(c.GetInt64("user_id"), req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
