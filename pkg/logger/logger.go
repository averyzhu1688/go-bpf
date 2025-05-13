package logger

import (
	"os"
	"path/filepath"
	"time"

	"bpf.com/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// global log
var Log *zap.Logger

// Init logger
func InitLogger() error {
	cfg := config.GetAppConfig().Log

	// set log level
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	//set console  encoder
	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	//set log ouput with color
	if cfg.ColorOutput {
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
		}
	} else {
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		consoleEncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
		}
	}

	//set file log encoder
	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	//create cores
	cores := []zapcore.Core{}

	// add console output
	var consoleEncoder zapcore.Encoder
	if cfg.Format == "json" {
		consoleEncoder = zapcore.NewJSONEncoder(consoleEncoderConfig)
	} else {
		consoleEncoder = zapcore.NewConsoleEncoder(consoleEncoderConfig)
	}
	cores = append(cores, zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	))

	if cfg.EnableFile {
		logDir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		hook := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

		cores = append(cores, zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(hook),
			level,
		))
	}

	core := zapcore.NewTee(cores...)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(Log)

	return nil
}

// Get logger
func GetLogger() *zap.Logger {
	return Log
}

// close log
func CloseLogger() {
	if Log != nil {
		Log.Sync()
	}
}

// use zap for gin logger
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetAppConfig()
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		if cfg.Log.EnableConsole {
			cost := time.Since(start)
			cost = cost / time.Millisecond
			Log.Info("HTTP request",
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("uri", path),
				zap.Duration("total time", cost),
			)
		}
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*os.PathError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if se.Err != nil && se.Err.Error() == "broken pipe" {
							brokenPipe = true
						}
					}
				}

				httpRequest := c.Request.Method + " " + c.Request.URL.String()
				if brokenPipe {
					Log.Error("disconnect",
						zap.Any("error", err),
						zap.String("request", httpRequest),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					Log.Error("recovery exception",
						zap.Any("error", err),
						zap.String("request", httpRequest),
						zap.Stack("trace"),
					)
				} else {
					Log.Error("recovery exception",
						zap.Any("error", err),
						zap.String("request", httpRequest),
					)
				}

				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
