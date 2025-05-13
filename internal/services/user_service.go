package services

import (
	"errors"

	"bpf.com/internal/models"
	"bpf.com/internal/repository"
)

// user service interface
type IUserService interface {
	GetUserById(userId uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	ListUsers(page, pageSize int, search string) ([]*models.User, int64, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userId uint64) error
	HasPermission(userId uint64, permission string) (bool, error)
}

// implements IUserService
type UserService struct {
	userRepo repository.IUserRepository
}

// Create UserService
func NewUserService() IUserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// Find user by Id
func (s *UserService) GetUserById(userId uint64) (*models.User, error) {
	return s.userRepo.FindById(userId)
}

// Find user by name
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.FindByUsername(username)
}

// Find user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// Find user list
func (s *UserService) ListUsers(page, pageSize int, search string) ([]*models.User, int64, error) {

	return s.userRepo.List(page, pageSize, search)
}

// Create User
func (s *UserService) CreateUser(user *models.User) error {
	existsUser1, _ := s.userRepo.FindByUsername(user.Username)
	if existsUser1 != nil {
		return errors.New("user does not exist")
	}
	existsUser2, _ := s.userRepo.FindByEmail(user.Email)
	if existsUser2 != nil {
		return errors.New("email already exist")
	}
	return s.userRepo.Create(user)
}

// Update user
func (s *UserService) UpdateUser(user *models.User) error {
	existingUser, err := s.userRepo.FindById(user.Id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user does not exist")
	}
	if user.Username != existingUser.Username {
		conflictUser, _ := s.userRepo.FindByUsername(user.Username)
		if conflictUser != nil && conflictUser.Id != user.Id {
			return errors.New("username already exist")
		}
	}
	if user.Email != existingUser.Email {
		conflictUser, _ := s.userRepo.FindByEmail(user.Email)
		if conflictUser != nil && conflictUser.Id != user.Id {
			return errors.New("email already exist")
		}
	}
	return s.userRepo.Update(user)
}

// Delete user
func (s *UserService) DeleteUser(userId uint64) error {
	existingUser, err := s.userRepo.FindById(userId)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user does not exist")
	}
	return s.userRepo.Delete(userId)
}

// check user permission
func (s *UserService) HasPermission(userId uint64, permission string) (bool, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("user does not exist")
	}
	return user.HasPermission(permission), nil
}
