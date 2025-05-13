// internal/models/role.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

const (
	PermUserView   = "user:view"
	PermUserCreate = "user:create"
	PermUserEdit   = "user:edit"
	PermUserDelete = "user:delete"

	PermContentView   = "content:view"
	PermContentCreate = "content:create"
	PermContentEdit   = "content:edit"
	PermContentDelete = "content:delete"

	PermSystemConfig = "system:config"
	PermSystemLog    = "system:log"
	PermSystemBackup = "system:backup"

	PermAll = "*"
)

type Permissions []string

func (p *Permissions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}
	return json.Unmarshal(bytes, p)
}

func (p Permissions) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

func (p Permissions) HasPermission(permission string) bool {
	for _, perm := range p {
		if perm == permission || perm == "*" {
			return true
		}
	}
	return false
}

func (p *Permissions) AddPermission(permission string) {
	for _, perm := range *p {
		if perm == permission {
			return
		}
	}
	*p = append(*p, permission)
}

func (p *Permissions) RemovePermission(permission string) {
	for i, perm := range *p {
		if perm == permission {
			*p = append((*p)[:i], (*p)[i+1:]...)
			return
		}
	}
}

type Role struct {
	BaseModel
	Name        string      `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Code        string      `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Description string      `gorm:"size:200" json:"description"`
	Permissions Permissions `gorm:"type:json" json:"permissions"`
}

func (Role) TableName() string {
	return "t_sys_roles"
}

func (r *Role) HasPermission(permission string) bool {
	return r.Permissions.HasPermission(permission)
}

func (r *Role) AddPermission(permission string) {
	r.Permissions.AddPermission(permission)
}

func (r *Role) RemovePermission(permission string) {
	r.Permissions.RemovePermission(permission)
}

var (
	AdminRole = &Role{
		Name:        "管理员",
		Code:        RoleAdmin,
		Description: "系统管理员，拥有所有权限",
		Permissions: Permissions{"*"},
	}

	UserRole = &Role{
		Name:        "普通用户",
		Code:        RoleUser,
		Description: "普通用户，拥有基本权限",
		Permissions: Permissions{
			"user:view",
			"user:edit",
			"content:view",
			"content:create",
			"content:edit",
		},
	}
	GuestRole = &Role{
		Name:        "访客",
		Code:        RoleGuest,
		Description: "访客，仅拥有查看权限",
		Permissions: Permissions{
			"user:view",
			"content:view",
		},
	}
)

func GetPredefinedRole(code string) *Role {
	switch code {
	case RoleAdmin:
		return AdminRole
	case RoleUser:
		return UserRole
	case RoleGuest:
		return GuestRole
	default:
		return nil
	}
}
