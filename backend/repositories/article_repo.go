package repositories

import (
	"time"

	"github.com/kubektl/v0-blog-backend/models"
	"gorm.io/gorm"
)

type ArticleFilter struct {
	Status   string
	Tag      string
	AuthorID int64
	TeamID   int64
	Query    string
	Page     int
	PerPage  int
	Since    *time.Time
	AuthorIDs []int64
}

type ArticleRepository interface {
	FindByID(id int64) (*models.Article, error)
	FindBySlug(slug string) (*models.Article, error)
	Create(article *models.Article) error
	Update(article *models.Article) error
	Delete(id int64) error
	List(filter ArticleFilter) ([]models.Article, int64, error)
	Trending(limit int, since time.Time) ([]models.Article, error)
	IncrementComments(id int64, delta int) error
	SetViewerFlags(articles []models.Article, userID int64) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) baseQuery() *gorm.DB {
	return r.db.Model(&models.Article{}).Preload("Author").Preload("Team")
}

func (r *articleRepository) applyFilter(q *gorm.DB, filter ArticleFilter) *gorm.DB {
	if filter.Status != "" {
		q = q.Where("articles.status = ?", filter.Status)
	}
	if filter.Tag != "" {
		q = q.Where("? = ANY(articles.tags)", filter.Tag)
	}
	if filter.AuthorID > 0 {
		q = q.Where("articles.author_id = ?", filter.AuthorID)
	}
	if filter.TeamID > 0 {
		q = q.Where("articles.team_id = ?", filter.TeamID)
	}
	if filter.Query != "" {
		pattern := "%" + filter.Query + "%"
		q = q.Where("articles.title ILIKE ? OR articles.excerpt ILIKE ?", pattern, pattern)
	}
	if filter.Since != nil {
		q = q.Where("articles.created_at >= ?", *filter.Since)
	}
	if len(filter.AuthorIDs) > 0 {
		q = q.Where("articles.author_id IN ?", filter.AuthorIDs)
	}
	return q
}

func (r *articleRepository) FindByID(id int64) (*models.Article, error) {
	var article models.Article
	err := r.baseQuery().First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := r.baseQuery().Where("slug = ?", slug).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) Delete(id int64) error {
	return r.db.Delete(&models.Article{}, id).Error
}

func (r *articleRepository) List(filter ArticleFilter) ([]models.Article, int64, error) {
	q := r.applyFilter(r.baseQuery(), filter)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 10
	}
	offset := (filter.Page - 1) * filter.PerPage

	var articles []models.Article
	err := q.Order("articles.created_at DESC").Limit(filter.PerPage).Offset(offset).Find(&articles).Error
	return articles, total, err
}

func (r *articleRepository) Trending(limit int, since time.Time) ([]models.Article, error) {
	var articles []models.Article
	err := r.baseQuery().
		Where("articles.status = ? AND articles.created_at >= ?", models.ArticleStatusPublished, since).
		Order("(articles.likes_count + articles.comments_count) DESC, articles.created_at DESC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}

func (r *articleRepository) IncrementComments(id int64, delta int) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).
		UpdateColumn("comments_count", gorm.Expr("GREATEST(comments_count + ?, 0)", delta)).Error
}

func (r *articleRepository) SetViewerFlags(articles []models.Article, userID int64) error {
	if userID == 0 || len(articles) == 0 {
		return nil
	}
	ids := make([]int64, len(articles))
	for i, a := range articles {
		ids[i] = a.ID
	}

	var likedIDs []int64
	r.db.Model(&models.Like{}).Where("user_id = ? AND article_id IN ?", userID, ids).Pluck("article_id", &likedIDs)
	likedSet := map[int64]bool{}
	for _, id := range likedIDs {
		likedSet[id] = true
	}

	var bookmarkedIDs []int64
	r.db.Model(&models.Bookmark{}).Where("user_id = ? AND article_id IN ?", userID, ids).Pluck("article_id", &bookmarkedIDs)
	bookmarkSet := map[int64]bool{}
	for _, id := range bookmarkedIDs {
		bookmarkSet[id] = true
	}

	for i := range articles {
		v := likedSet[articles[i].ID]
		articles[i].LikedByMe = &v
		b := bookmarkSet[articles[i].ID]
		articles[i].BookmarkedByMe = &b
	}
	return nil
}
