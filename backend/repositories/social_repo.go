package repositories

import (
	"github.com/kubektl/v0-blog-backend/models"
	"gorm.io/gorm"
)

type SocialRepository interface {
	ToggleLike(articleID, userID int64) (liked bool, err error)
	ToggleBookmark(articleID, userID int64) (bookmarked bool, err error)
	ToggleFollow(followerID, followingID int64) (following bool, err error)
	GetBookmarks(userID int64, page, perPage int) ([]models.Article, int64, error)
	GetComments(articleID int64, page, perPage int) ([]models.Comment, int64, error)
	CreateComment(comment *models.Comment) error
	DeleteComment(commentID, userID int64) (articleID int64, err error)
	GetFollowingIDs(userID int64) ([]int64, error)
}

type socialRepository struct {
	db *gorm.DB
}

func NewSocialRepository(db *gorm.DB) SocialRepository {
	return &socialRepository{db: db}
}

func (r *socialRepository) ToggleLike(articleID, userID int64) (bool, error) {
	var liked bool
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Like
		err := tx.Where("article_id = ? AND user_id = ?", articleID, userID).First(&existing).Error
		if err == nil {
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Article{}).Where("id = ?", articleID).
				UpdateColumn("likes_count", gorm.Expr("GREATEST(likes_count - 1, 0)")).Error; err != nil {
				return err
			}
			liked = false
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := tx.Create(&models.Like{ArticleID: articleID, UserID: userID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Article{}).Where("id = ?", articleID).
			UpdateColumn("likes_count", gorm.Expr("likes_count + 1")).Error; err != nil {
			return err
		}
		liked = true
		return nil
	})
	return liked, err
}

func (r *socialRepository) ToggleBookmark(articleID, userID int64) (bool, error) {
	var bookmarked bool
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Bookmark
		err := tx.Where("article_id = ? AND user_id = ?", articleID, userID).First(&existing).Error
		if err == nil {
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
			bookmarked = false
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := tx.Create(&models.Bookmark{ArticleID: articleID, UserID: userID}).Error; err != nil {
			return err
		}
		bookmarked = true
		return nil
	})
	return bookmarked, err
}

func (r *socialRepository) ToggleFollow(followerID, followingID int64) (bool, error) {
	var following bool
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Follow
		err := tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existing).Error
		if err == nil {
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.User{}).Where("id = ?", followingID).
				UpdateColumn("followers", gorm.Expr("GREATEST(followers - 1, 0)")).Error; err != nil {
				return err
			}
			following = false
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := tx.Create(&models.Follow{FollowerID: followerID, FollowingID: followingID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).Where("id = ?", followingID).
			UpdateColumn("followers", gorm.Expr("followers + 1")).Error; err != nil {
			return err
		}
		following = true
		return nil
	})
	return following, err
}

func (r *socialRepository) GetBookmarks(userID int64, page, perPage int) ([]models.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	var total int64
	r.db.Model(&models.Bookmark{}).Where("user_id = ?", userID).Count(&total)

	var articles []models.Article
	err := r.db.Model(&models.Article{}).
		Joins("JOIN bookmarks ON bookmarks.article_id = articles.id").
		Preload("Author").
		Where("bookmarks.user_id = ?", userID).
		Order("bookmarks.created_at DESC").
		Limit(perPage).Offset(offset).
		Find(&articles).Error
	return articles, total, err
}

func (r *socialRepository) GetComments(articleID int64, page, perPage int) ([]models.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var total int64
	r.db.Model(&models.Comment{}).Where("article_id = ?", articleID).Count(&total)

	var comments []models.Comment
	err := r.db.Preload("User").
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Limit(perPage).Offset(offset).
		Find(&comments).Error
	return comments, total, err
}

func (r *socialRepository) CreateComment(comment *models.Comment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		return tx.Model(&models.Article{}).Where("id = ?", comment.ArticleID).
			UpdateColumn("comments_count", gorm.Expr("comments_count + 1")).Error
	})
}

func (r *socialRepository) DeleteComment(commentID, userID int64) (int64, error) {
	var articleID int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var comment models.Comment
		if err := tx.First(&comment, commentID).Error; err != nil {
			return err
		}
		if comment.UserID != userID {
			return gorm.ErrInvalidData
		}
		articleID = comment.ArticleID
		if err := tx.Delete(&comment).Error; err != nil {
			return err
		}
		return tx.Model(&models.Article{}).Where("id = ?", articleID).
			UpdateColumn("comments_count", gorm.Expr("GREATEST(comments_count - 1, 0)")).Error
	})
	return articleID, err
}

func (r *socialRepository) GetFollowingIDs(userID int64) ([]int64, error) {
	var ids []int64
	err := r.db.Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Pluck("following_id", &ids).Error
	return ids, err
}
