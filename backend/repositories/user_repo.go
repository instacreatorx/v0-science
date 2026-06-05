package repositories

import (
	"github.com/kubektl/v0-blog-backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id int64) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	IncrementFollowers(userID int64, delta int) error
	IsFollowing(followerID, followingID int64) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) IncrementFollowers(userID int64, delta int) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("followers", gorm.Expr("GREATEST(followers + ?, 0)", delta)).Error
}

func (r *userRepository) IsFollowing(followerID, followingID int64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}
