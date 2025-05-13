// internal/models/user.go
package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	RoleSuperuser = "superuser"
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleGuest     = "guest"
)

type User struct {
	BaseModel
	Username  string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string     `gorm:"size:100;not null" json:"-"`
	Email     string     `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone     string     `gorm:"size:20" json:"phone"`
	Nickname  string     `gorm:"size:50" json:"nickname"`
	RoleId    uint       `gorm:"default:3" json:"role_id"`
	Role      *Role      `gorm:"foreignKey:RoleId" json:"role,omitempty"`
	Status    int        `gorm:"default:1" json:"status"`
	LastLogin *time.Time `json:"last_login"`
}

func (User) TableName() string {
	return "t_sys_users"
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("密码不能为空")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsActive() bool {
	return u.Status == 1
}

func (u *User) IsAdmin() bool {
	if u.Role != nil {
		return u.Role.Code == RoleAdmin
	}
	return false
}

func (u *User) HasPermission(permission string) bool {
	if u.Role != nil {
		return u.Role.HasPermission(permission)
	}
	return false
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
}
