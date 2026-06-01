package handlers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/kubektl/v0-blog-backend/database"
	"github.com/kubektl/v0-blog-backend/models"
)

func calculateReadTime(content string) string {
	words := len(strings.Fields(content))
	minutes := int(math.Ceil(float64(words) / 200.0))
	if minutes < 1 {
		minutes = 1
	}
	return strconv.Itoa(minutes) + " min read"
}

func GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	tag := c.Query("tag")
	authorID := c.Query("author_id")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 50 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	// Build query
	query := `
		SELECT a.id, a.title, a.excerpt, a.content, a.image, a.author_id, a.tags, 
		       a.read_time, a.is_member_only, a.likes_count, a.comments_count, a.created_at, a.updated_at,
		       u.id, u.name, u.avatar, u.bio, u.followers
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM articles WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	if tag != "" {
		argCount++
		query += " AND $" + strconv.Itoa(argCount) + " = ANY(a.tags)"
		countQuery += " AND $" + strconv.Itoa(argCount) + " = ANY(tags)"
		args = append(args, tag)
	}

	if authorID != "" {
		argCount++
		query += " AND a.author_id = $" + strconv.Itoa(argCount)
		countQuery += " AND author_id = $" + strconv.Itoa(argCount)
		args = append(args, authorID)
	}

	// Get total count
	var total int64
	err := database.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count articles"})
		return
	}

	// Add pagination
	query += " ORDER BY a.created_at DESC LIMIT $" + strconv.Itoa(argCount+1) + " OFFSET $" + strconv.Itoa(argCount+2)
	args = append(args, perPage, offset)

	rows, err := database.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		var a models.Article
		var author models.User
		err := rows.Scan(
			&a.ID, &a.Title, &a.Excerpt, &a.Content, &a.Image, &a.AuthorID, &a.Tags,
			&a.ReadTime, &a.IsMemberOnly, &a.LikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt,
			&author.ID, &author.Name, &author.Avatar, &author.Bio, &author.Followers,
		)
		if err != nil {
			continue
		}
		a.Author = &author
		a.Content = "" // Don't send full content in list
		articles = append(articles, a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       articles,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetArticle(c *gin.Context) {
	articleID := c.Param("id")

	var a models.Article
	var author models.User
	err := database.DB.QueryRow(context.Background(), `
		SELECT a.id, a.title, a.excerpt, a.content, a.image, a.author_id, a.tags, 
		       a.read_time, a.is_member_only, a.likes_count, a.comments_count, a.created_at, a.updated_at,
		       u.id, u.name, u.avatar, u.bio, u.followers
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE a.id = $1
	`, articleID).Scan(
		&a.ID, &a.Title, &a.Excerpt, &a.Content, &a.Image, &a.AuthorID, &a.Tags,
		&a.ReadTime, &a.IsMemberOnly, &a.LikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt,
		&author.ID, &author.Name, &author.Avatar, &author.Bio, &author.Followers,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch article"})
		}
		return
	}

	a.Author = &author
	c.JSON(http.StatusOK, a)
}

func CreateArticle(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req models.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	readTime := calculateReadTime(req.Content)

	var a models.Article
	err := database.DB.QueryRow(context.Background(), `
		INSERT INTO articles (title, excerpt, content, image, author_id, tags, read_time, is_member_only)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, excerpt, content, image, author_id, tags, read_time, is_member_only, likes_count, comments_count, created_at, updated_at
	`, req.Title, req.Excerpt, req.Content, req.Image, userID, req.Tags, readTime, req.IsMemberOnly,
	).Scan(&a.ID, &a.Title, &a.Excerpt, &a.Content, &a.Image, &a.AuthorID, &a.Tags, &a.ReadTime, &a.IsMemberOnly, &a.LikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	// Get author info
	var author models.User
	database.DB.QueryRow(context.Background(),
		"SELECT id, name, avatar, bio, followers FROM users WHERE id = $1", userID,
	).Scan(&author.ID, &author.Name, &author.Avatar, &author.Bio, &author.Followers)
	a.Author = &author

	c.JSON(http.StatusCreated, a)
}

func UpdateArticle(c *gin.Context) {
	userID := c.GetInt64("user_id")
	articleID := c.Param("id")

	// Check ownership
	var ownerID int64
	err := database.DB.QueryRow(context.Background(),
		"SELECT author_id FROM articles WHERE id = $1", articleID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit this article"})
		return
	}

	var req models.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	readTime := calculateReadTime(req.Content)

	var a models.Article
	err = database.DB.QueryRow(context.Background(), `
		UPDATE articles SET
			title = COALESCE(NULLIF($1, ''), title),
			excerpt = COALESCE(NULLIF($2, ''), excerpt),
			content = COALESCE(NULLIF($3, ''), content),
			image = COALESCE(NULLIF($4, ''), image),
			tags = COALESCE($5, tags),
			read_time = $6,
			is_member_only = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING id, title, excerpt, content, image, author_id, tags, read_time, is_member_only, likes_count, comments_count, created_at, updated_at
	`, req.Title, req.Excerpt, req.Content, req.Image, req.Tags, readTime, req.IsMemberOnly, articleID,
	).Scan(&a.ID, &a.Title, &a.Excerpt, &a.Content, &a.Image, &a.AuthorID, &a.Tags, &a.ReadTime, &a.IsMemberOnly, &a.LikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, a)
}

func DeleteArticle(c *gin.Context) {
	userID := c.GetInt64("user_id")
	articleID := c.Param("id")

	// Check ownership
	var ownerID int64
	err := database.DB.QueryRow(context.Background(),
		"SELECT author_id FROM articles WHERE id = $1", articleID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this article"})
		return
	}

	_, err = database.DB.Exec(context.Background(), "DELETE FROM articles WHERE id = $1", articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func LikeArticle(c *gin.Context) {
	userID := c.GetInt64("user_id")
	articleID := c.Param("id")

	// Check if already liked
	var exists bool
	database.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM likes WHERE article_id = $1 AND user_id = $2)",
		articleID, userID).Scan(&exists)

	if exists {
		// Unlike
		_, err := database.DB.Exec(context.Background(),
			"DELETE FROM likes WHERE article_id = $1 AND user_id = $2", articleID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike article"})
			return
		}
		database.DB.Exec(context.Background(),
			"UPDATE articles SET likes_count = likes_count - 1 WHERE id = $1", articleID)
		c.JSON(http.StatusOK, gin.H{"liked": false, "message": "Article unliked"})
	} else {
		// Like
		_, err := database.DB.Exec(context.Background(),
			"INSERT INTO likes (article_id, user_id) VALUES ($1, $2)", articleID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like article"})
			return
		}
		database.DB.Exec(context.Background(),
			"UPDATE articles SET likes_count = likes_count + 1 WHERE id = $1", articleID)
		c.JSON(http.StatusOK, gin.H{"liked": true, "message": "Article liked"})
	}
}

func BookmarkArticle(c *gin.Context) {
	userID := c.GetInt64("user_id")
	articleID := c.Param("id")

	// Check if already bookmarked
	var exists bool
	database.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM bookmarks WHERE article_id = $1 AND user_id = $2)",
		articleID, userID).Scan(&exists)

	if exists {
		// Remove bookmark
		_, err := database.DB.Exec(context.Background(),
			"DELETE FROM bookmarks WHERE article_id = $1 AND user_id = $2", articleID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove bookmark"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"bookmarked": false, "message": "Bookmark removed"})
	} else {
		// Add bookmark
		_, err := database.DB.Exec(context.Background(),
			"INSERT INTO bookmarks (article_id, user_id) VALUES ($1, $2)", articleID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bookmark article"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"bookmarked": true, "message": "Article bookmarked"})
	}
}

func GetBookmarks(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage

	var total int64
	database.DB.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM bookmarks WHERE user_id = $1", userID).Scan(&total)

	rows, err := database.DB.Query(context.Background(), `
		SELECT a.id, a.title, a.excerpt, a.image, a.author_id, a.tags, 
		       a.read_time, a.is_member_only, a.likes_count, a.comments_count, a.created_at, a.updated_at,
		       u.id, u.name, u.avatar
		FROM bookmarks b
		JOIN articles a ON b.article_id = a.id
		JOIN users u ON a.author_id = u.id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookmarks"})
		return
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		var a models.Article
		var author models.User
		rows.Scan(
			&a.ID, &a.Title, &a.Excerpt, &a.Image, &a.AuthorID, &a.Tags,
			&a.ReadTime, &a.IsMemberOnly, &a.LikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt,
			&author.ID, &author.Name, &author.Avatar,
		)
		a.Author = &author
		articles = append(articles, a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       articles,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetComments(c *gin.Context) {
	articleID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage

	var total int64
	database.DB.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM comments WHERE article_id = $1", articleID).Scan(&total)

	rows, err := database.DB.Query(context.Background(), `
		SELECT c.id, c.article_id, c.user_id, c.content, c.created_at, c.updated_at,
		       u.id, u.name, u.avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.article_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`, articleID, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		var user models.User
		rows.Scan(
			&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
			&user.ID, &user.Name, &user.Avatar,
		)
		comment.User = &user
		comments = append(comments, comment)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       comments,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func CreateComment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	articleID := c.Param("id")

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check article exists
	var exists bool
	database.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM articles WHERE id = $1)", articleID).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	var comment models.Comment
	err := database.DB.QueryRow(context.Background(), `
		INSERT INTO comments (article_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, article_id, user_id, content, created_at, updated_at
	`, articleID, userID, req.Content,
	).Scan(&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Update comments count
	database.DB.Exec(context.Background(),
		"UPDATE articles SET comments_count = comments_count + 1 WHERE id = $1", articleID)

	// Get user info
	var user models.User
	database.DB.QueryRow(context.Background(),
		"SELECT id, name, avatar FROM users WHERE id = $1", userID,
	).Scan(&user.ID, &user.Name, &user.Avatar)
	comment.User = &user

	c.JSON(http.StatusCreated, comment)
}

func DeleteComment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	commentID := c.Param("commentId")

	// Check ownership
	var ownerID int64
	var articleID int64
	err := database.DB.QueryRow(context.Background(),
		"SELECT user_id, article_id FROM comments WHERE id = $1", commentID).Scan(&ownerID, &articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this comment"})
		return
	}

	_, err = database.DB.Exec(context.Background(), "DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	// Update comments count
	database.DB.Exec(context.Background(),
		"UPDATE articles SET comments_count = comments_count - 1 WHERE id = $1", articleID)

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func FollowUser(c *gin.Context) {
	followerID := c.GetInt64("user_id")
	followingID := c.Param("id")

	if strconv.FormatInt(followerID, 10) == followingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot follow yourself"})
		return
	}

	// Check if already following
	var exists bool
	database.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)",
		followerID, followingID).Scan(&exists)

	if exists {
		// Unfollow
		_, err := database.DB.Exec(context.Background(),
			"DELETE FROM follows WHERE follower_id = $1 AND following_id = $2", followerID, followingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
			return
		}
		database.DB.Exec(context.Background(),
			"UPDATE users SET followers = followers - 1 WHERE id = $1", followingID)
		c.JSON(http.StatusOK, gin.H{"following": false, "message": "Unfollowed user"})
	} else {
		// Follow
		_, err := database.DB.Exec(context.Background(),
			"INSERT INTO follows (follower_id, following_id) VALUES ($1, $2)", followerID, followingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
			return
		}
		database.DB.Exec(context.Background(),
			"UPDATE users SET followers = followers + 1 WHERE id = $1", followingID)
		c.JSON(http.StatusOK, gin.H{"following": true, "message": "Following user"})
	}
}

func GetUserArticles(c *gin.Context) {
	userID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage

	var total int64
	database.DB.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM articles WHERE author_id = $1", userID).Scan(&total)

	rows, err := database.DB.Query(context.Background(), `
		SELECT a.id, a.title, a.excerpt, a.image, a.author_id, a.tags, 
		       a.read_time, a.is_member_only, a.likes_count, a.comments_count, a.created_at, a.updated_at,
		       u.id, u.name, u.avatar
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE a.author_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		var a models.Article
		var author models.User
		rows.Scan(
			&a.ID, &a.Title, &a.Excerpt, &a.Image, &a.AuthorID, &a.Tags,
			&a.ReadTime, &a.IsMemberOnly, &a.LikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt,
			&author.ID, &author.Name, &author.Avatar,
		)
		a.Author = &author
		articles = append(articles, a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       articles,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func SearchArticles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage

	searchPattern := "%" + query + "%"

	var total int64
	database.DB.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM articles WHERE title ILIKE $1 OR excerpt ILIKE $1", searchPattern).Scan(&total)

	rows, err := database.DB.Query(context.Background(), `
		SELECT a.id, a.title, a.excerpt, a.image, a.author_id, a.tags, 
		       a.read_time, a.is_member_only, a.likes_count, a.comments_count, a.created_at, a.updated_at,
		       u.id, u.name, u.avatar
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE a.title ILIKE $1 OR a.excerpt ILIKE $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`, searchPattern, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search articles"})
		return
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		var a models.Article
		var author models.User
		rows.Scan(
			&a.ID, &a.Title, &a.Excerpt, &a.Image, &a.AuthorID, &a.Tags,
			&a.ReadTime, &a.IsMemberOnly, &a.LikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt,
			&author.ID, &author.Name, &author.Avatar,
		)
		a.Author = &author
		articles = append(articles, a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       articles,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetTrendingArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "6"))
	if limit < 1 || limit > 20 {
		limit = 6
	}

	// Get articles from last 7 days with most engagement
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	rows, err := database.DB.Query(context.Background(), `
		SELECT a.id, a.title, a.excerpt, a.image, a.author_id, a.tags, 
		       a.read_time, a.is_member_only, a.likes_count, a.comments_count, a.created_at, a.updated_at,
		       u.id, u.name, u.avatar
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE a.created_at >= $1
		ORDER BY (a.likes_count + a.comments_count) DESC, a.created_at DESC
		LIMIT $2
	`, sevenDaysAgo, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trending articles"})
		return
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		var a models.Article
		var author models.User
		rows.Scan(
			&a.ID, &a.Title, &a.Excerpt, &a.Image, &a.AuthorID, &a.Tags,
			&a.ReadTime, &a.IsMemberOnly, &a.LikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt,
			&author.ID, &author.Name, &author.Avatar,
		)
		a.Author = &author
		articles = append(articles, a)
	}

	c.JSON(http.StatusOK, articles)
}
