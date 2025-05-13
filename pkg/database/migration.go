// pkg/database/migration.go
package database

import (
	"bpf.com/internal/models"
	"bpf.com/pkg/logger"
	"go.uber.org/zap"
)

func RunMigrations() error {
	logger.GetLogger().Info("开始执行数据库迁移...")
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
	); err != nil {
		logger.GetLogger().Error("数据库迁移失败", zap.Error(err))
		return err
	}

	logger.GetLogger().Info("数据库迁移完成")
	return nil
}

func InitAdminUser() error {
	if err := initRoles(); err != nil {
		return err
	}
	if err := initAdmin(); err != nil {
		return err
	}
	return nil
}

func initRoles() error {
	logger.GetLogger().Info("检查并初始化角色...")
	var count int64
	if err := DB.Model(&models.Role{}).Count(&count).Error; err != nil {
		logger.GetLogger().Error("查询角色表失败", zap.Error(err))
		return err
	}

	if count == 0 {
		logger.GetLogger().Info("创建默认角色...")
		superuserRole := &models.Role{
			Name:        "超级管理员",
			Code:        models.RoleSuperuser,
			Description: "系统最高管理员，拥有所有权限",
			Permissions: models.Permissions{models.PermAll},
		}

		adminRole := &models.Role{
			Name:        "管理员",
			Code:        models.RoleAdmin,
			Description: "系统管理员，拥有所有权限",
			Permissions: models.Permissions{models.PermAll},
		}

		userRole := &models.Role{
			Name:        "普通用户",
			Code:        models.RoleUser,
			Description: "普通用户，拥有基本权限",
			Permissions: models.Permissions{
				models.PermUserView,
				models.PermUserEdit,
				models.PermContentView,
				models.PermContentCreate,
				models.PermContentEdit,
			},
		}

		guestRole := &models.Role{
			Name:        "访客",
			Code:        models.RoleGuest,
			Description: "访客，仅拥有查看权限",
			Permissions: models.Permissions{
				models.PermUserView,
				models.PermContentView,
			},
		}

		roles := []*models.Role{superuserRole, adminRole, userRole, guestRole}
		if err := DB.Create(&roles).Error; err != nil {
			logger.GetLogger().Error("创建默认角色失败", zap.Error(err))
			return err
		}

		logger.GetLogger().Info("默认角色创建成功")
	}
	return nil
}

func initAdmin() error {
	logger.GetLogger().Info("检查并初始化管理员账户...")
	var count int64
	if err := DB.Model(&models.User{}).Where("role_id = ?", 1).Count(&count).Error; err != nil {
		logger.GetLogger().Error("查询管理员账户失败", zap.Error(err))
		return err
	}

	if count == 0 {
		logger.GetLogger().Info("创建默认管理员账户...")
		admin := &models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Nickname: "系统管理员",
			RoleId:   1, // 管理员角色ID为1
			Status:   models.StatusActive,
		}
		if err := admin.SetPassword("123456"); err != nil {
			logger.GetLogger().Error("设置管理员密码失败", zap.Error(err))
			return err
		}
		if err := DB.Create(admin).Error; err != nil {
			logger.GetLogger().Error("创建管理员账户失败", zap.Error(err))
			return err
		}

		logger.GetLogger().Info("默认管理员账户创建成功")
	}

	return nil
}
