package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// service degradation
func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				log.Printf("[PANIC RECOVER] error: %v\trace: %s", err, string(stack))
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "server internal error",
				})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
