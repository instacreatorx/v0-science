package services

import (
	"time"

	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/database"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/repositories"
	"github.com/kubektl/v0-blog-backend/services/state"
	"gorm.io/gorm"
)

type TeamService struct {
	teams    repositories.TeamRepository
	users    repositories.UserRepository
	articles repositories.ArticleRepository
}

func NewTeamService(teams repositories.TeamRepository, users repositories.UserRepository, articles repositories.ArticleRepository) *TeamService {
	return &TeamService{teams: teams, users: users, articles: articles}
}

func (s *TeamService) Create(ownerID int64, req models.CreateTeamRequest) (*models.Team, error) {
	slug, err := database.UniqueTeamSlug(req.Slug, 0)
	if err != nil {
		return nil, err
	}

	team := &models.Team{
		Name:    req.Name,
		Slug:    slug,
		Bio:     req.Bio,
		Avatar:  req.Avatar,
		OwnerID: ownerID,
	}
	if team.Avatar == "" {
		team.Avatar = "https://api.dicebear.com/7.x/shapes/svg?seed=" + slug
	}

	if err := s.teams.Create(team); err != nil {
		return nil, err
	}
	if err := s.teams.AddMember(&models.TeamMember{
		TeamID: team.ID,
		UserID: ownerID,
		Role:   models.TeamRoleOwner,
	}); err != nil {
		return nil, err
	}
	return s.teams.FindByID(team.ID)
}

func (s *TeamService) GetBySlug(slug string) (*models.Team, error) {
	team, err := s.teams.FindBySlug(slug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return team, nil
}

func (s *TeamService) Update(userID, teamID int64, req models.UpdateTeamRequest) (*models.Team, error) {
	if err := s.requireRole(teamID, userID, models.TeamRoleOwner, models.TeamRoleAdmin); err != nil {
		return nil, err
	}
	team, err := s.teams.FindByID(teamID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Bio != "" {
		team.Bio = req.Bio
	}
	if req.Avatar != "" {
		team.Avatar = req.Avatar
	}
	if err := s.teams.Update(team); err != nil {
		return nil, err
	}
	return s.teams.FindByID(teamID)
}

func (s *TeamService) ListMembers(teamID int64) ([]models.TeamMember, error) {
	return s.teams.GetMembers(teamID)
}

func (s *TeamService) AddMember(actorID, teamID int64, req models.AddTeamMemberRequest) (*models.TeamMember, error) {
	if err := s.requireRole(teamID, actorID, models.TeamRoleOwner, models.TeamRoleAdmin); err != nil {
		return nil, err
	}
	if _, err := s.users.FindByID(req.UserID); err != nil {
		return nil, apperrors.ErrNotFound
	}
	validRoles := map[string]bool{
		models.TeamRoleAdmin: true, models.TeamRoleEditor: true, models.TeamRoleWriter: true,
	}
	if !validRoles[req.Role] {
		return nil, apperrors.ErrBadRequest
	}
	member := &models.TeamMember{TeamID: teamID, UserID: req.UserID, Role: req.Role}
	if err := s.teams.AddMember(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *TeamService) RemoveMember(actorID, teamID, targetUserID int64) error {
	if err := s.requireRole(teamID, actorID, models.TeamRoleOwner, models.TeamRoleAdmin); err != nil {
		return err
	}
	team, err := s.teams.FindByID(teamID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	if team.OwnerID == targetUserID {
		return apperrors.ErrForbidden
	}
	return s.teams.RemoveMember(teamID, targetUserID)
}

func (s *TeamService) ListUserTeams(userID int64) ([]models.Team, error) {
	return s.teams.ListUserTeams(userID)
}

func (s *TeamService) SubmitVerification(userID, teamID int64, req models.TeamVerifyRequest) (*models.TeamVerificationRequest, error) {
	if err := s.requireRole(teamID, userID, models.TeamRoleOwner, models.TeamRoleAdmin); err != nil {
		return nil, err
	}
	team, err := s.teams.FindByID(teamID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	if team.VerifiedAt != nil {
		return nil, apperrors.ErrConflict
	}
	if _, err := s.teams.FindPendingVerification(teamID); err == nil {
		return nil, apperrors.ErrConflict
	}

	vr := &models.TeamVerificationRequest{
		TeamID:      teamID,
		SubmittedBy: userID,
		ProofText:   req.ProofText,
		ProofURL:    req.ProofURL,
		Status:      models.VerificationStatusPending,
	}
	if err := s.teams.CreateVerificationRequest(vr); err != nil {
		return nil, err
	}
	return vr, nil
}

func (s *TeamService) ListVerificationRequests(status string, page, perPage int) ([]models.TeamVerificationRequest, int64, int, error) {
	reqs, total, err := s.teams.ListVerificationRequests(status, page, perPage)
	if err != nil {
		return nil, 0, 0, err
	}
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	return reqs, total, totalPages, nil
}

func (s *TeamService) ApproveVerification(adminID, requestID int64) (*models.Team, error) {
	req, team, err := s.loadVerificationRequest(requestID)
	if err != nil {
		return nil, err
	}
	if err := state.ValidateVerificationTransition(req.Status, models.VerificationStatusApproved); err != nil {
		return nil, err
	}
	now := time.Now()
	req.Status = models.VerificationStatusApproved
	req.ReviewedBy = &adminID
	req.ReviewedAt = &now
	team.VerifiedAt = &now
	if err := s.teams.UpdateVerificationRequest(req); err != nil {
		return nil, err
	}
	if err := s.teams.Update(team); err != nil {
		return nil, err
	}
	return s.teams.FindByID(team.ID)
}

func (s *TeamService) RejectVerification(adminID, requestID int64, reason string) (*models.TeamVerificationRequest, error) {
	req, _, err := s.loadVerificationRequest(requestID)
	if err != nil {
		return nil, err
	}
	if err := state.ValidateVerificationTransition(req.Status, models.VerificationStatusRejected); err != nil {
		return nil, err
	}
	now := time.Now()
	req.Status = models.VerificationStatusRejected
	req.ReviewedBy = &adminID
	req.ReviewedAt = &now
	req.RejectionReason = reason
	if err := s.teams.UpdateVerificationRequest(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *TeamService) loadVerificationRequest(requestID int64) (*models.TeamVerificationRequest, *models.Team, error) {
	var req models.TeamVerificationRequest
	if err := database.DB.First(&req, requestID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, apperrors.ErrNotFound
		}
		return nil, nil, err
	}
	team, err := s.teams.FindByID(req.TeamID)
	if err != nil {
		return nil, nil, apperrors.ErrNotFound
	}
	return &req, team, nil
}

func (s *TeamService) requireRole(teamID, userID int64, roles ...string) error {
	member, err := s.teams.GetMember(teamID, userID)
	if err != nil {
		return apperrors.ErrForbidden
	}
	for _, r := range roles {
		if member.Role == r {
			return nil
		}
	}
	return apperrors.ErrForbidden
}
