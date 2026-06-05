package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/middleware"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/services"
)

type AdminHandler struct {
	teams *services.TeamService
}

func NewAdminHandler(teams *services.TeamService) *AdminHandler {
	return &AdminHandler{teams: teams}
}

func (h *AdminHandler) ListVerificationRequests(c *gin.Context) {
	page, perPage := pagination(c)
	status := c.DefaultQuery("status", models.VerificationStatusPending)
	reqs, total, totalPages, err := h.teams.ListVerificationRequests(status, page, perPage)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: reqs, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}

func (h *AdminHandler) ApproveVerification(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid request ID")
		return
	}
	team, err := h.teams.ApproveVerification(c.GetInt64("user_id"), requestID)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, team)
}

func (h *AdminHandler) RejectVerification(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid request ID")
		return
	}
	var req models.RejectVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	vr, err := h.teams.RejectVerification(c.GetInt64("user_id"), requestID, req.Reason)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, vr)
}
