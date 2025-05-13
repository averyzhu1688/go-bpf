package repository

import (
	"bpf.com/internal/models"
	"bpf.com/pkg/database"
	"gorm.io/gorm"
)

// User repository interface
type IUserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint64) error
	FindById(id uint64) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	List(page, size int, query string) ([]*models.User, int64, error)
}

// UserRepository implements IUserRepository
type UserRepository struct {
	db *gorm.DB
}

// create UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// save user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// update user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// delete user
func (r *UserRepository) Delete(id uint64) error {
	return r.db.Delete(&models.User{}, id).Error
}

// find user by id
func (r *UserRepository) FindById(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// find user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// find user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// find user list
func (r *UserRepository) List(page, size int, query string) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	db := r.db.Model(&models.User{}).Preload("Role")

	//add query params
	if query != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	//count
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	//page query
	offset := (page - 1) * size
	err = db.Offset(offset).Limit(size).Limit(size).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
