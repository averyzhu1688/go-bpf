package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"bpf.com/pkg/config"
	"github.com/gin-gonic/gin"
)

const (
	Reset      = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
)

func statusCodeColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return BoldGreen
	case code >= 300 && code < 400:
		return BoldBlue
	case code >= 400 && code < 500:
		return BoldYellow
	default:
		return BoldRed
	}
}

func methodColor(method string) string {
	switch method {
	case "GET":
		return Blue
	case "POST":
		return Green
	case "PUT":
		return Yellow
	case "DELETE":
		return Red
	case "PATCH":
		return Cyan
	case "HEAD":
		return Magenta
	case "OPTIONS":
		return White
	default:
		return Reset
	}
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := config.GetAppConfig()
		startTime := time.Now()

		var requestBody []byte
		if ctx.Request.Body != nil {
			requestBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		ctx.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		if cfg.Server.EnableRequestLog {
			statusColor := statusCodeColor(statusCode)
			methodColor := methodColor(reqMethod)
			latencyStr := fmt.Sprintf("%.3fms", float64(latencyTime.Microseconds())/1000.0)
			fmt.Printf("[%s][GIN] %s%-7s%s | %s%3d%s | %13s | %15s | %s\n",
				time.Now().Format("2006-01-02 15:04:05"),
				methodColor, reqMethod, Reset,
				statusColor, statusCode, Reset,
				latencyStr,
				clientIP,
				reqUri,
			)

			if statusCode >= 400 {
				fmt.Printf("%s[ERROR]%s Request: %s\n", BoldRed, Reset, string(requestBody))
				fmt.Printf("%s[ERROR]%s Response: %s\n", BoldRed, Reset, blw.body.String())
			}
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
