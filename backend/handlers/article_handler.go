package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/middleware"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/repositories"
	"github.com/kubektl/v0-blog-backend/services"
)

type ArticleHandler struct {
	articles *services.ArticleService
}

func NewArticleHandler(articles *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{articles: articles}
}

func pagination(c *gin.Context) (page, perPage int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ = strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 50 {
		perPage = 10
	}
	return page, perPage
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	page, perPage := pagination(c)
	filter := repositories.ArticleFilter{
		Status:  models.ArticleStatusPublished,
		Tag:     c.Query("tag"),
		Page:    page,
		PerPage: perPage,
	}
	if authorID := c.Query("author_id"); authorID != "" {
		id, err := strconv.ParseInt(authorID, 10, 64)
		if err == nil {
			filter.AuthorID = id
		}
	}
	articles, total, totalPages, err := h.articles.List(filter, middleware.GetOptionalUserID(c))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	articles = h.articles.StripContentForList(articles)
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: articles, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}

func (h *ArticleHandler) GetMyArticles(c *gin.Context) {
	page, perPage := pagination(c)
	status := c.Query("status")
	articles, total, totalPages, err := h.articles.MyArticles(c.GetInt64("user_id"), status, page, perPage)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: articles, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}

func (h *ArticleHandler) GetFeed(c *gin.Context) {
	page, perPage := pagination(c)
	articles, total, totalPages, err := h.articles.Feed(c.GetInt64("user_id"), page, perPage)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	articles = h.articles.StripContentForList(articles)
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: articles, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}

func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Search query required")
		return
	}
	page, perPage := pagination(c)
	filter := repositories.ArticleFilter{
		Status:  models.ArticleStatusPublished,
		Query:   query,
		Page:    page,
		PerPage: perPage,
	}
	articles, total, totalPages, err := h.articles.List(filter, middleware.GetOptionalUserID(c))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	articles = h.articles.StripContentForList(articles)
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: articles, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}

func (h *ArticleHandler) GetTrendingArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "6"))
	if limit < 1 || limit > 20 {
		limit = 6
	}
	articles, err := h.articles.Trending(limit, middleware.GetOptionalUserID(c))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	article, err := h.articles.GetBySlug(c.Param("slug"), middleware.GetOptionalUserID(c))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid article ID")
		return
	}
	article, err := h.articles.GetByID(id, middleware.GetOptionalUserID(c))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req models.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	article, err := h.articles.Create(c.GetInt64("user_id"), req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid article ID")
		return
	}
	var req models.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	article, err := h.articles.Update(c.GetInt64("user_id"), id, req)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid article ID")
		return
	}
	if err := h.articles.Delete(c.GetInt64("user_id"), id); err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	h.transitionArticle(c, models.ArticleStatusPublished)
}

func (h *ArticleHandler) UnpublishArticle(c *gin.Context) {
	h.transitionArticle(c, models.ArticleStatusDraft)
}

func (h *ArticleHandler) ArchiveArticle(c *gin.Context) {
	h.transitionArticle(c, models.ArticleStatusArchived)
}

func (h *ArticleHandler) transitionArticle(c *gin.Context, status string) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid article ID")
		return
	}
	article, err := h.articles.Transition(c.GetInt64("user_id"), id, status)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, article)
}
