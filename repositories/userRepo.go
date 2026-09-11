package repositories

import (
	"fmt"
	"linkedveil/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindById(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Delete(id uint) error
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) Delete(id uint) error {
	return u.db.Delete(&models.User{}, id).Error
}

func (u *userRepository) FindById(id uint) (*models.User, error) {
	var user models.User

	err := u.db.Find(&user, id).Error
	if err != nil { 
		return nil, fmt.Errorf("couldn't find user by id: %w", err)
	}

	return &user, nil
}

func (u *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User

	err := u.db.Where("username = ?", username).Find(&user).Error
	if err != nil { 
		return nil, fmt.Errorf("couldn't find user by username: %w", err)
	}

	return &user, nil
}
