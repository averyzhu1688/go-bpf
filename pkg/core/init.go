package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bpf.com/pkg/cache"
	"bpf.com/pkg/config"
	"bpf.com/pkg/database"
	"bpf.com/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	server *http.Server
	engine *gin.Engine
}

// Init Config
func InitConfig(configPath string) error {
	return config.InitConfig(configPath)
}

// Init Logger
func InitLogger() error {
	return logger.InitLogger()
}

// Init Database
func InitDatabase() error {
	return database.InitDatabase()
}

// Init Cache
func InitCache() error {
	return cache.InitRedisCache()
}

// Init gin engine
func InitGin() *gin.Engine {
	gin.SetMode(config.GetAppConfig().Server.Mode)
	if config.GetAppConfig().Server.DisableDebug {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}
	router := gin.New()
	router.Use(logger.GinLogger())
	router.Use(logger.GinRecovery(true))
	return router
}

// create Application
func NewApplication(router *gin.Engine) *Application {
	serverConfig := config.GetAppConfig().Server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler:      router,
		ReadTimeout:  serverConfig.ReadTimeout * time.Second,
		WriteTimeout: serverConfig.WriteTimeout * time.Second,
	}
	return &Application{
		server: srv,
		engine: router,
	}
}

// get Gin engine
func (app *Application) Engine() *gin.Engine {
	return app.engine
}

// start server
func (app *Application) Run() {
	//start http server
	go func() {
		host := config.GetAppConfig().Server.Host
		port := config.GetAppConfig().Server.Port

		var serverURL string
		if host == "0.0.0.0" {
			serverURL = fmt.Sprintf("http://localhost:%d", port)
		} else {
			serverURL = fmt.Sprintf("http://%s:%d", host, port)
		}
		logger.GetLogger().Info(fmt.Sprintf("server is start....[%s]", serverURL))
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Fatal("server start fail:", zap.Error(err))
		}
	}()
	//close server
	app.gracefulShutdown()
}

// Secure shutdown service
func (app *Application) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.GetLogger().Info("server close...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		logger.GetLogger().Fatal("server closed fail",
			zap.Error(err))
	}
	database.CloseDatabase()
	logger.CloseLogger()
	logger.GetLogger().Info("server closed")
}
