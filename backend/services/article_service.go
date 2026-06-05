package services

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/database"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/repositories"
	"github.com/kubektl/v0-blog-backend/services/state"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ArticleService struct {
	articles repositories.ArticleRepository
	social   repositories.SocialRepository
}

func NewArticleService(articles repositories.ArticleRepository, social repositories.SocialRepository) *ArticleService {
	return &ArticleService{articles: articles, social: social}
}

func CalculateReadTime(content string) string {
	words := len(strings.Fields(content))
	minutes := int(math.Ceil(float64(words) / 200.0))
	if minutes < 1 {
		minutes = 1
	}
	return strconv.Itoa(minutes) + " min read"
}

func (s *ArticleService) List(filter repositories.ArticleFilter, viewerID int64) ([]models.Article, int64, int, error) {
	if filter.Status == "" && filter.AuthorID == 0 {
		filter.Status = models.ArticleStatusPublished
	}
	articles, total, err := s.articles.List(filter)
	if err != nil {
		return nil, 0, 0, err
	}
	_ = s.articles.SetViewerFlags(articles, viewerID)
	totalPages := int(math.Ceil(float64(total) / float64(filter.PerPage)))
	return articles, total, totalPages, nil
}

func (s *ArticleService) Trending(limit int, viewerID int64) ([]models.Article, error) {
	since := time.Now().AddDate(0, 0, -7)
	articles, err := s.articles.Trending(limit, since)
	if err != nil {
		return nil, err
	}
	_ = s.articles.SetViewerFlags(articles, viewerID)
	return articles, nil
}

func (s *ArticleService) Feed(userID int64, page, perPage int) ([]models.Article, int64, int, error) {
	followingIDs, err := s.social.GetFollowingIDs(userID)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(followingIDs) == 0 {
		return []models.Article{}, 0, 0, nil
	}
	filter := repositories.ArticleFilter{
		Status:    models.ArticleStatusPublished,
		AuthorIDs: followingIDs,
		Page:      page,
		PerPage:   perPage,
	}
	return s.List(filter, userID)
}

func (s *ArticleService) GetByID(id int64, viewerID int64) (*models.Article, error) {
	article, err := s.articles.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return s.applyAccess(article, viewerID)
}

func (s *ArticleService) GetBySlug(slug string, viewerID int64) (*models.Article, error) {
	article, err := s.articles.FindBySlug(slug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return s.applyAccess(article, viewerID)
}

func (s *ArticleService) applyAccess(article *models.Article, viewerID int64) (*models.Article, error) {
	if article.Status != models.ArticleStatusPublished {
		if viewerID != article.AuthorID {
			return nil, apperrors.ErrNotFound
		}
	}
	if article.IsMemberOnly && viewerID == 0 {
		return nil, apperrors.ErrMemberOnly
	}
	articles := []models.Article{*article}
	_ = s.articles.SetViewerFlags(articles, viewerID)
	*article = articles[0]
	return article, nil
}

func (s *ArticleService) Create(userID int64, req models.CreateArticleRequest) (*models.Article, error) {
	slug, err := database.UniqueSlug(req.Title, 0)
	if err != nil {
		return nil, err
	}

	excerpt := req.Excerpt
	if excerpt == "" && req.Content != "" {
		words := strings.Fields(req.Content)
		if len(words) > 30 {
			excerpt = strings.Join(words[:30], " ") + "..."
		} else {
			excerpt = req.Content
		}
	}

	article := &models.Article{
		Title:        req.Title,
		Slug:         slug,
		Excerpt:      excerpt,
		Content:      req.Content,
		Image:        req.Image,
		AuthorID:     userID,
		TeamID:       req.TeamID,
		Tags:         pq.StringArray(req.Tags),
		ReadTime:     CalculateReadTime(req.Content),
		Status:       models.ArticleStatusDraft,
		IsMemberOnly: req.IsMemberOnly,
	}
	if err := s.articles.Create(article); err != nil {
		return nil, err
	}
	return s.articles.FindByID(article.ID)
}

func (s *ArticleService) Update(userID, articleID int64, req models.UpdateArticleRequest) (*models.Article, error) {
	article, err := s.articles.FindByID(articleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	if article.AuthorID != userID {
		return nil, apperrors.ErrForbidden
	}

	if req.Title != "" {
		article.Title = req.Title
		if slug, err := database.UniqueSlug(req.Title, article.ID); err == nil {
			article.Slug = slug
		}
	}
	if req.Excerpt != "" {
		article.Excerpt = req.Excerpt
	}
	if req.Content != "" {
		article.Content = req.Content
		article.ReadTime = CalculateReadTime(req.Content)
	}
	if req.Image != "" {
		article.Image = req.Image
	}
	if req.Tags != nil {
		article.Tags = pq.StringArray(req.Tags)
	}
	if req.TeamID != nil {
		article.TeamID = req.TeamID
	}
	if req.IsMemberOnly != nil {
		article.IsMemberOnly = *req.IsMemberOnly
	}

	if err := s.articles.Update(article); err != nil {
		return nil, err
	}
	return s.articles.FindByID(article.ID)
}

func (s *ArticleService) Delete(userID, articleID int64) error {
	article, err := s.articles.FindByID(articleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrNotFound
		}
		return err
	}
	if article.AuthorID != userID {
		return apperrors.ErrForbidden
	}
	return s.articles.Delete(articleID)
}

func (s *ArticleService) Transition(userID, articleID int64, toStatus string) (*models.Article, error) {
	article, err := s.articles.FindByID(articleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	if article.AuthorID != userID {
		return nil, apperrors.ErrForbidden
	}
	if err := state.ValidateArticleTransition(article.Status, toStatus); err != nil {
		return nil, err
	}

	article.Status = toStatus
	if toStatus == models.ArticleStatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}
	if toStatus == models.ArticleStatusDraft {
		article.PublishedAt = nil
	}

	if err := s.articles.Update(article); err != nil {
		return nil, err
	}
	return s.articles.FindByID(article.ID)
}

func (s *ArticleService) MyArticles(userID int64, status string, page, perPage int) ([]models.Article, int64, int, error) {
	filter := repositories.ArticleFilter{
		AuthorID: userID,
		Status:   status,
		Page:     page,
		PerPage:  perPage,
	}
	articles, total, err := s.articles.List(filter)
	if err != nil {
		return nil, 0, 0, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return articles, total, totalPages, nil
}

func (s *ArticleService) StripContentForList(articles []models.Article) []models.Article {
	for i := range articles {
		articles[i].Content = ""
	}
	return articles
}
