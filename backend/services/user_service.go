package services

import (
	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/repositories"
	"gorm.io/gorm"
)

type UserService struct {
	users repositories.UserRepository
}

func NewUserService(users repositories.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetByID(id int64, viewerID int64) (*models.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	if viewerID > 0 && viewerID != id {
		following, _ := s.users.IsFollowing(viewerID, id)
		user.IsFollowing = &following
	}
	return user, nil
}

func (s *UserService) GetCurrent(id int64) (*models.User, error) {
	return s.GetByID(id, 0)
}

func (s *UserService) Update(id int64, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	if req.Locale != "" && req.Locale != "en" && req.Locale != "fa" {
		return nil, apperrors.ErrBadRequest
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Locale != "" {
		user.Locale = req.Locale
	}

	if err := s.users.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}
