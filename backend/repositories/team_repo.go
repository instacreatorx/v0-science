package repositories

import (
	"github.com/kubektl/v0-blog-backend/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(team *models.Team) error
	FindByID(id int64) (*models.Team, error)
	FindBySlug(slug string) (*models.Team, error)
	Update(team *models.Team) error
	AddMember(member *models.TeamMember) error
	RemoveMember(teamID, userID int64) error
	GetMembers(teamID int64) ([]models.TeamMember, error)
	GetMember(teamID, userID int64) (*models.TeamMember, error)
	CreateVerificationRequest(req *models.TeamVerificationRequest) error
	FindPendingVerification(teamID int64) (*models.TeamVerificationRequest, error)
	ListVerificationRequests(status string, page, perPage int) ([]models.TeamVerificationRequest, int64, error)
	UpdateVerificationRequest(req *models.TeamVerificationRequest) error
	ListUserTeams(userID int64) ([]models.Team, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepository) FindByID(id int64) (*models.Team, error) {
	var team models.Team
	err := r.db.Preload("Owner").First(&team, id).Error
	if err != nil {
		return nil, err
	}
	team.IsVerified = team.VerifiedAt != nil
	return &team, nil
}

func (r *teamRepository) FindBySlug(slug string) (*models.Team, error) {
	var team models.Team
	err := r.db.Preload("Owner").Where("slug = ?", slug).First(&team).Error
	if err != nil {
		return nil, err
	}
	team.IsVerified = team.VerifiedAt != nil
	return &team, nil
}

func (r *teamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

func (r *teamRepository) AddMember(member *models.TeamMember) error {
	return r.db.Create(member).Error
}

func (r *teamRepository) RemoveMember(teamID, userID int64) error {
	return r.db.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&models.TeamMember{}).Error
}

func (r *teamRepository) GetMembers(teamID int64) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.db.Preload("User").Where("team_id = ?", teamID).Find(&members).Error
	return members, err
}

func (r *teamRepository) GetMember(teamID, userID int64) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *teamRepository) CreateVerificationRequest(req *models.TeamVerificationRequest) error {
	return r.db.Create(req).Error
}

func (r *teamRepository) FindPendingVerification(teamID int64) (*models.TeamVerificationRequest, error) {
	var req models.TeamVerificationRequest
	err := r.db.Where("team_id = ? AND status = ?", teamID, models.VerificationStatusPending).
		Order("created_at DESC").First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *teamRepository) ListVerificationRequests(status string, page, perPage int) ([]models.TeamVerificationRequest, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	q := r.db.Model(&models.TeamVerificationRequest{}).Preload("Team")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	q.Count(&total)

	var reqs []models.TeamVerificationRequest
	err := q.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&reqs).Error
	return reqs, total, err
}

func (r *teamRepository) UpdateVerificationRequest(req *models.TeamVerificationRequest) error {
	return r.db.Save(req).Error
}

func (r *teamRepository) ListUserTeams(userID int64) ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userID).
		Preload("Owner").
		Find(&teams).Error
	for i := range teams {
		teams[i].IsVerified = teams[i].VerifiedAt != nil
	}
	return teams, err
}
