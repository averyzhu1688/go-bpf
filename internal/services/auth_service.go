package services

import (
	"errors"

	"bpf.com/internal/models"
	"bpf.com/internal/repository"
	"bpf.com/pkg/utils"
	"gorm.io/gorm"
)

// user auth interface
type IAuthService interface {
	Register(username, password, email, nickname string) (*models.User, error)
	Login(username, password string) (string, string, *models.User, error)
	RefreshToken(refreshtoken string) (string, error)
	VerifyToken(token string) (*models.User, error)
	ChangePassword(userId uint64, oldPassword, newPassword string) error
}

// auth implements
type AuthService struct {
	userRepo repository.IUserRepository
}

// create new AuthService
func NewAuthService() IAuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Register
func (s *AuthService) Register(username, password, email, nickname string) (*models.User, error) {
	//check user exists
	querUser1, _ := s.userRepo.FindByUsername(username)
	if querUser1 != nil {
		return nil, errors.New("username already exists")
	}
	//check email
	querUser2, _ := s.userRepo.FindByEmail(email)
	if querUser2 != nil {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Nickname: nickname,
		RoleId:   2,
		Status:   models.StatusActive,
	}
	//set password
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	//save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login
func (s *AuthService) Login(username, password string) (string, string, *models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, errors.New("user does not exist")
		}
		return "", "", nil, err
	}
	//check user status
	if !user.IsActive() {
		return "", "", nil, errors.New("user disabled")
	}
	//check password
	if !user.CheckPassword(password) {
		return "", "", nil, errors.New("password error")
	}

	//update lastedLogin time
	user.UpdateLastLogin()
	if err := s.userRepo.Update(user); err != nil {
		return "", "", nil, err
	}

	//create accessToken and refreshToken
	accessToken, err := utils.GenerateAccessToken(user.Id)
	if err != nil {
		return "", "", nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.Id)
	if err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, user, nil
}

// refresh token
func (s *AuthService) RefreshToken(refreshtoken string) (string, error) {
	claims, err := utils.ParseRefreshToken(refreshtoken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}
	//check user
	user, err := s.userRepo.FindById(claims.UserId)
	if err != nil {
		return "", errors.New("user does not exist")
	}
	if !user.IsActive() {
		return "", errors.New("user disabled")
	}
	//generate new token
	accessToken, err := utils.GenerateAccessToken(user.Id)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// verify token
func (s *AuthService) VerifyToken(token string) (*models.User, error) {
	claims, err := utils.ParseAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	user, err := s.userRepo.FindById(claims.UserId)
	if err != nil {
		return nil, errors.New("user does not exist")
	}
	if !user.IsActive() {
		return nil, errors.New("user disabled")
	}
	return user, nil
}

// change password
func (s *AuthService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindById(userID)
	if err != nil {
		return errors.New("user does not exist")
	}
	if user.CheckPassword(oldPassword) {
		return errors.New("old password error")
	}
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}
	return s.userRepo.Update(user)
}
