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

type TeamHandler struct {
	teams *services.TeamService
}

func NewTeamHandler(teams *services.TeamService) *TeamHandler {
	return &TeamHandler{teams: teams}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req models.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	team, err := h.teams.Create(c.GetInt64("user_id"), req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) GetTeamBySlug(c *gin.Context) {
	team, err := h.teams.GetBySlug(c.Param("slug"))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid team ID")
		return
	}
	var req models.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	team, err := h.teams.Update(c.GetInt64("user_id"), teamID, req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) ListMembers(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid team ID")
		return
	}
	members, err := h.teams.ListMembers(teamID)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *TeamHandler) AddMember(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid team ID")
		return
	}
	var req models.AddTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	member, err := h.teams.AddMember(c.GetInt64("user_id"), teamID, req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, member)
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid team ID")
		return
	}
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid user ID")
		return
	}
	if err := h.teams.RemoveMember(c.GetInt64("user_id"), teamID, userID); err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

func (h *TeamHandler) ListMyTeams(c *gin.Context) {
	teams, err := h.teams.ListUserTeams(c.GetInt64("user_id"))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) SubmitVerification(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid team ID")
		return
	}
	var req models.TeamVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	vr, err := h.teams.SubmitVerification(c.GetInt64("user_id"), teamID, req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, vr)
}
