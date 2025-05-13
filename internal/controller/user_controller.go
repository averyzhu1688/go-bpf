package controller

import (
	"context"
	"strconv"
	"time"

	"bpf.com/internal/models"
	"bpf.com/internal/services"
	"bpf.com/pkg/cache"
	"bpf.com/pkg/utils"
	"github.com/gin-gonic/gin"
)

// User Controller
type UserController struct {
	userService services.IUserService
}

// Create UserController
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Create User Request
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
	RoleId   uint   `json:"role_id" binding:"required"`
}

// Update User Request
type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	RoleId   uint   `json:"roleId" binding:"required"`
}

// Query User Request
type QueryUserRequest struct {
	Search   string `json:"search"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
}

// Get User list
func (c *UserController) GetUsers(ctx *gin.Context) {
	//get page pramas
	//page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	//pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	//search := ctx.DefaultQuery("search", "")

	var req QueryUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, err.Error(), nil)
		return
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	Search := req.Search

	users, total, err := c.userService.ListUsers(pageNum, pageSize, Search)
	if err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}

	//Use Redis Cache
	cache.GetGlobalCache().Set(context.Background(), "user-views-info", users, 10*time.Minute)

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":       user.Id,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"role": gin.H{
				"id":   user.Role.Id,
				"name": user.Role.Name,
				"code": user.Role.Code,
			},
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}
	utils.Success(ctx, gin.H{
		"list":  userList,
		"total": total,
		"page":  pageNum,
		"size":  pageSize,
	})
}

// Find user by Id
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, "无效的用户ID", nil)
		return
	}
	user, err := c.userService.GetUserById(id)
	if err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}
	if user == nil {
		utils.FailWithMessage(ctx, utils.NOT_FOUND, "user not found", nil)
		return
	}

	utils.Success(ctx, gin.H{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
		"role": gin.H{
			"id":   user.Role.Id,
			"name": user.Role.Name,
			"code": user.Role.Code,
		},
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// Create User
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, err.Error(), nil)
		return
	}
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Nickname: req.Nickname,
		RoleId:   req.RoleId,
	}
	if err := user.SetPassword(req.Password); err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, "set password fail: "+err.Error(), nil)
		return
	}
	if err := c.userService.CreateUser(user); err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}
	utils.Success(ctx, gin.H{
		"user_id":  user.Id,
		"username": user.Username,
	})
}

// Update User
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, "invalid user", nil)
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, err.Error(), nil)
		return
	}

	user, err := c.userService.GetUserById(userId)
	if err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}
	if user == nil {
		utils.FailWithMessage(ctx, utils.NOT_FOUND, "user not found", nil)
		return
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.RoleId = req.RoleId

	if err := c.userService.UpdateUser(user); err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}

	utils.SuccessWithMessage(ctx, "update user successfully", nil)
}

// Delete user
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.FailWithMessage(ctx, utils.INVALID_PARAMS, "invalid userId", nil)
		return
	}
	if err := c.userService.DeleteUser(userId); err != nil {
		utils.FailWithMessage(ctx, utils.ERROR, err.Error(), nil)
		return
	}
	utils.SuccessWithMessage(ctx, "delete user successfully", nil)
}
