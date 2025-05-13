package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"bpf.com/internal/services"
	"bpf.com/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Auth middleware
func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "authorization token not found",
			})
			ctx.Abort()
			return
		}

		// Check Bearer Prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "auth error 'Bearer {token}'",
			})
			ctx.Abort()
			return
		}
		token := parts[1]
		claims, err := utils.ParseAccessToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid token: " + err.Error(),
			})
			ctx.Abort()
			return
		}
		// set user to ctx
		ctx.Set("userId", claims.UserId)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}

// Role auth middleware
func RoleAuth(roleCode string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorization",
			})
			ctx.Abort()
			return
		}
		userService := services.NewUserService()
		user, err := userService.GetUserById(userId.(uint64))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "get user fail: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		if user == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "user not found",
			})
			ctx.Abort()
			return
		}
		fmt.Println("User Role:", user.Role.Code, "Required Role:", roleCode)
		if user.Role == nil || user.Role.Code != roleCode {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "access control require role: " + roleCode + " 角色",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// Permission Auth middleware
func PermissionAuth(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthenticated user",
			})
			ctx.Abort()
			return
		}
		userService := services.NewUserService()
		hasPermission, err := userService.HasPermission(userId.(uint64), permission)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "permission check fail: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		if !hasPermission {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "access control require permission: " + permission,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// need role and permission
func RoleAndPermissionAuth(roleCode string, permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthen",
			})
			ctx.Abort()
			return
		}

		userService := services.NewUserService()
		user, err := userService.GetUserById(userId.(uint64))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "get user fail: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		if user == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "user not found",
			})
			ctx.Abort()
			return
		}

		if user.Role == nil || user.Role.Code != roleCode {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "access control require role: " + roleCode,
			})
			ctx.Abort()
			return
		}

		hasPermission, err := userService.HasPermission(userId.(uint64), permission)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "check permission fail: " + err.Error(),
			})
			ctx.Abort()
			return
		}
		if !hasPermission {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "access control require permission: " + permission,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// Role or Permession
func RoleOrPermissionAuth(roleCode string, permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorization",
			})
			ctx.Abort()
			return
		}
		userService := services.NewUserService()
		user, err := userService.GetUserById(userId.(uint64))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "get user info fail: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		if user == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "user not found",
			})
			ctx.Abort()
			return
		}

		hasRole := user.Role != nil && user.Role.Code == roleCode
		hasPermission, err := userService.HasPermission(userId.(uint64), permission)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "check permission fail: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		if !hasRole && !hasPermission {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "access control require role: " + roleCode + " or permission:" + permission,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// At least one authorization is required
func AnyPermissionAuth(permissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorization",
			})
			ctx.Abort()
			return
		}

		userService := services.NewUserService()
		for _, permission := range permissions {
			hasPermission, err := userService.HasPermission(userId.(uint64), permission)
			if err != nil {
				continue
			}
			if hasPermission {
				ctx.Next()
				return
			}
		}

		permissionList := strings.Join(permissions, ", ")
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "permission control requires one of the following authorizations: " + permissionList,
		})
		ctx.Abort()
	}
}

func AllPermissionsAuth(permissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorization",
			})
			ctx.Abort()
			return
		}

		userService := services.NewUserService()
		for _, permission := range permissions {
			hasPermission, err := userService.HasPermission(userId.(uint64), permission)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "check permission fail: " + err.Error(),
				})
				ctx.Abort()
				return
			}
			if !hasPermission {

				ctx.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "access control require permission : " + permission,
				})
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}
