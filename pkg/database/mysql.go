package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"bpf.com/pkg/config"
	"bpf.com/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init database
func InitDatabase() error {
	cfg := config.GetAppConfig().Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	var logWriter io.Writer
	if cfg.DisableSqlLog {
		logWriter = io.Discard
	} else {
		logWriter = os.Stdout
	}
	var logLevel gormlogger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = gormlogger.Silent
	case "error":
		logLevel = gormlogger.Error
	case "warn":
		logLevel = gormlogger.Warn
	case "info":
		logLevel = gormlogger.Info
	default:
		logLevel = gormlogger.Info
	}

	gormLogger := gormlogger.New(
		log.New(logWriter, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             cfg.SlowThreshold * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("connector db fail: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("get connector pool fail: %w", err)
	}
	//set conn pool
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Second)

	// test conn
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("test connector : %w", err)
	}

	logger.GetLogger().Info("db connector successfully")

	// move database
	if cfg.AutoMigrate {
		if err := RunMigrations(); err != nil {
			logger.GetLogger().Error("database move fail", zap.Error(err))
			return err
		}
	}

	// init admin
	if cfg.InitAdmin {
		if err := InitAdminUser(); err != nil {
			logger.GetLogger().Error("init admin fail", zap.Error(err))
			return err
		}
	}
	return nil
}

// close
func CloseDatabase() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			logger.GetLogger().Error("get connector fail", zap.Error(err))
			return
		}
		if err := sqlDB.Close(); err != nil {
			logger.GetLogger().Error("close connector fail", zap.Error(err))
			return
		}
		logger.GetLogger().Info("db connector is close")
	}
}

// get db
func GetDB() *gorm.DB {
	return DB
}
