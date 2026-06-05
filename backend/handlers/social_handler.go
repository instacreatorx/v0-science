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

type SocialHandler struct {
	social *services.SocialService
}

func NewSocialHandler(social *services.SocialService) *SocialHandler {
	return &SocialHandler{social: social}
}

func (h *SocialHandler) LikeArticle(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid article ID")
		return
	}
	liked, err := h.social.ToggleLike(articleID, c.GetInt64("user_id"))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	msg := "Article liked"
	if !liked {
		msg = "Article unliked"
	}
	c.JSON(http.StatusOK, gin.H{"liked": liked, "message": msg})
}

func (h *SocialHandler) BookmarkArticle(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid article ID")
		return
	}
	bookmarked, err := h.social.ToggleBookmark(articleID, c.GetInt64("user_id"))
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	msg := "Article bookmarked"
	if !bookmarked {
		msg = "Bookmark removed"
	}
	c.JSON(http.StatusOK, gin.H{"bookmarked": bookmarked, "message": msg})
}

func (h *SocialHandler) GetBookmarks(c *gin.Context) {
	page, perPage := pagination(c)
	articles, total, totalPages, err := h.social.GetBookmarks(c.GetInt64("user_id"), page, perPage)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: articles, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}

func (h *SocialHandler) GetComments(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid article ID")
		return
	}
	page, perPage := pagination(c)
	comments, total, totalPages, err := h.social.GetComments(articleID, page, perPage)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: comments, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}

func (h *SocialHandler) CreateComment(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid article ID")
		return
	}
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, err.Error())
		return
	}
	comment, err := h.social.CreateComment(articleID, c.GetInt64("user_id"), req.Content)
	if err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *SocialHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("commentId"), 10, 64)
	if err != nil {
		middleware.RespondError(c, apperrors.ErrBadRequest, "Invalid comment ID")
		return
	}
	if err := h.social.DeleteComment(commentID, c.GetInt64("user_id")); err != nil {
		middleware.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
