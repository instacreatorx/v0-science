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

type UserHandler struct {
	users    *services.UserService
	social   *services.SocialService
	articles *services.ArticleService
}

func NewUserHandler(users *services.UserService, social *services.SocialService, articles *services.ArticleService) *UserHandler {
	return &UserHandler{users: users, social: social, articles: articles}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid user ID")
		return
	}
	viewerID := middleware.GetOptionalUserID(c)
	user, err := h.users.GetByID(id, viewerID)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) FollowUser(c *gin.Context) {
	followingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid user ID")
		return
	}
	following, err := h.social.ToggleFollow(c.GetInt64("user_id"), followingID)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	msg := "Following user"
	if !following {
		msg = "Unfollowed user"
	}
	c.JSON(http.StatusOK, gin.H{"following": following, "message": msg})
}

func (h *UserHandler) GetUserArticles(c *gin.Context) {
	authorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid user ID")
		return
	}
	page, perPage := pagination(c)
	articles, total, totalPages, err := h.articles.List(repositories.ArticleFilter{
		AuthorID: authorID,
		Status:   models.ArticleStatusPublished,
		Page:     page,
		PerPage:  perPage,
	}, middleware.GetOptionalUserID(c))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	articles = h.articles.StripContentForList(articles)
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: articles, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}
