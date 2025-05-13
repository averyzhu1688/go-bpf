package controller

import (
	"bpf.com/internal/services"
	"bpf.com/pkg/utils"
	"github.com/gin-gonic/gin"
)

// auth Controller
type AuthController struct {
	authService services.IAuthService
}

// Create AuthController
func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// Register quest params
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
}

// Login request params
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Refresh Token Request params
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Change password request params
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

// User Register
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, err.Error(), nil)
		return
	}
	user, err := c.authService.Register(req.Username, req.Password, req.Email, req.Nickname)
	if err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}
	utils.Success(ctx, gin.H{
		"user_id":  user.Id,
		"username": user.Username,
	})
}

// Login
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, err.Error(), nil)
		return
	}
	accessToken, refreshToken, user, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}
	utils.Success(ctx, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":       user.Id,
			"username": user.Username,
			"nickname": user.Nickname,
			"email":    user.Email,
			"role": gin.H{
				"id":   user.Role.Id,
				"name": user.Role.Name,
				"code": user.Role.Code,
			},
		},
	})
}

// Refresh Token
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, err.Error(), nil)
		return
	}
	accessToken, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}
	utils.Success(ctx, gin.H{
		"access_token": accessToken,
	})
}

// Get User info
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	//get userId from ctx
	userId, exists := ctx.Get("userId")
	if !exists {
		utils.FailWithMessage(ctx, utils.UNAUTHORIZED, "user not found", nil)
		return
	}
	token := ctx.GetHeader("Authorization")
	if token == "" {
		utils.FailWithMessage(ctx, utils.UNAUTHORIZED, "token not found", nil)
		return
	}
	//remove Bearer prifix
	token = token[7:]
	user, err := c.authService.VerifyToken(token)
	if err != nil {
		utils.FailWithMessage(ctx, utils.UNAUTHORIZED, err.Error(), nil)
		return
	}
	//check userId
	if user.Id != userId.(uint64) {
		utils.FailWithMessage(ctx, utils.UNAUTHORIZED, "user authization fail", nil)
		return
	}
	utils.Success(ctx, gin.H{
		"user": gin.H{
			"id":       user.Id,
			"username": user.Username,
			"nickname": user.Nickname,
			"email":    user.Email,
			"role": gin.H{
				"id":   user.Role.Id,
				"name": user.Role.Name,
				"code": user.Role.Code,
			},
		},
	})
}

// Change Password
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, err.Error(), nil)
		return
	}
	userId, exists := ctx.Get("userId")
	if !exists {
		utils.FailWithMessage(ctx, utils.UNAUTHORIZED, "user not found", nil)
		return
	}
	err := c.authService.ChangePassword(userId.(uint64), req.OldPassword, req.NewPassword)
	if err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}
	utils.SuccessWithMessage(ctx, "change password successfully", nil)
}
