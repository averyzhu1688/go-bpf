package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

var globalConfig AppConfig

// app config
type AppConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
	Cache    CacheConfig
}

// server config
type ServerConfig struct {
	Host             string        `mapstructure:"host"`
	Port             int           `mapstructure:"port"`
	Mode             string        `mapstructure:"mode"`
	ReadTimeout      time.Duration `mapstructure:"readTimeout"`
	WriteTimeout     time.Duration `mapstructure:"writeTimeout"`
	DisableDebug     bool          `mapstructure:"disableDebug"`
	EnableRequestLog bool          `mapstructure:"enableRequestLog"`
}

// db config
type DatabaseConfig struct {
	Type            string        `mapstructure:"type"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
	LogLevel        string        `mapstructure:"logLevel"`
	SlowThreshold   time.Duration `mapstructure:"slowThreshold"`
	DisableSqlLog   bool          `mapstructure:"disableSqlLog"`
	AutoMigrate     bool          `mapstructure:"autoMigrate"`
	InitAdmin       bool          `mapstructure:"initAdmin"`
}

// jwt config
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	AccessTokenExp   time.Duration `mapstructure:"accessTokenExp"`
	RefreshTokenExp  time.Duration `mapstructure:"refreshTokenExp"`
	TokenIssuer      string        `mapstructure:"tokenIssuer"`
	RefreshTokenSize int           `mapstructure:"refreshTokenSize"`
}

// cache config
type CacheConfig struct {
	Type         string        `mapstructure:"type"`
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"poolSize"`
	MinIdleConns int           `mapstructure:"minIdleConns"`
	MaxRetries   int           `mapstructure:"maxRetries"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	DefaultTTL   time.Duration `mapstructure:"defaultTTL"`
	Prefix       string        `mapstructure:"prefix"`
	EnableLog    bool          `mapstructure:"enableLog"`
}

// log config
type LogConfig struct {
	Level         string `mapstructure:"level"`
	Filename      string `mapstructure:"filename"`
	MaxSize       int    `mapstructure:"maxSize"`
	MaxBackups    int    `mapstructure:"maxBackups"`
	MaxAge        int    `mapstructure:"maxAge"`
	Compress      bool   `mapstructure:"compress"`
	EnableFile    bool   `mapstructure:"enableFile"`
	Format        string `mapstructure:"format"`
	ColorOutput   bool   `mapstructure:"colorOutput"`
	EnableConsole bool   `mapstructure:"enableConsole"`
}

// Init all config
func InitConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	//read file...
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read the config fail:%w", err)
	}

	//parse the config parames to AppConfig struct
	if err := viper.Unmarshal(&globalConfig); err != nil {
		return fmt.Errorf("parse the config fail: %w", err)
	}
	log.Println("config file process successfully!")
	return nil
}

// get config
func GetAppConfig() *AppConfig {
	return &globalConfig
}
