package services

import (
	"math"

	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/repositories"
	"gorm.io/gorm"
)

type SocialService struct {
	social   repositories.SocialRepository
	articles repositories.ArticleRepository
}

func NewSocialService(social repositories.SocialRepository, articles repositories.ArticleRepository) *SocialService {
	return &SocialService{social: social, articles: articles}
}

func (s *SocialService) ToggleLike(articleID, userID int64) (bool, error) {
	return s.social.ToggleLike(articleID, userID)
}

func (s *SocialService) ToggleBookmark(articleID, userID int64) (bool, error) {
	return s.social.ToggleBookmark(articleID, userID)
}

func (s *SocialService) ToggleFollow(followerID, followingID int64) (bool, error) {
	if followerID == followingID {
		return false, apperrors.ErrBadRequest
	}
	return s.social.ToggleFollow(followerID, followingID)
}

func (s *SocialService) GetBookmarks(userID int64, page, perPage int) ([]models.Article, int64, int, error) {
	articles, total, err := s.social.GetBookmarks(userID, page, perPage)
	if err != nil {
		return nil, 0, 0, err
	}
	_ = s.articles.SetViewerFlags(articles, userID)
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return articles, total, totalPages, nil
}

func (s *SocialService) GetComments(articleID int64, page, perPage int) ([]models.Comment, int64, int, error) {
	comments, total, err := s.social.GetComments(articleID, page, perPage)
	if err != nil {
		return nil, 0, 0, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return comments, total, totalPages, nil
}

func (s *SocialService) CreateComment(articleID, userID int64, content string) (*models.Comment, error) {
	if _, err := s.articles.FindByID(articleID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	comment := &models.Comment{ArticleID: articleID, UserID: userID, Content: content}
	if err := s.social.CreateComment(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *SocialService) DeleteComment(commentID, userID int64) error {
	_, err := s.social.DeleteComment(commentID, userID)
	if err == gorm.ErrInvalidData {
		return apperrors.ErrForbidden
	}
	if err == gorm.ErrRecordNotFound {
		return apperrors.ErrNotFound
	}
	return err
}
