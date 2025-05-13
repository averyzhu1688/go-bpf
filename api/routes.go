package api

import (
	"time"

	"bpf.com/internal/controller"
	"bpf.com/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// Set system routes
func SetupRoutes(router *gin.Engine) {
	//add global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.Cors())
	//180 calls per minute
	router.Use(middleware.RateLimit(180, time.Minute))

	apiGroup := router.Group("/api/v1")
	{
		setupAuthRoutes(apiGroup)
		setupUserRoutes(apiGroup)
	}

}

// Auth routes
func setupAuthRoutes(apiGroup *gin.RouterGroup) {
	authController := controller.NewAuthController()
	//public route
	publicGroup := apiGroup.Group("/auth")
	{
		publicGroup.POST("/register", authController.Register)
		publicGroup.POST("/login", authController.Login)
		publicGroup.POST("/refresh", authController.RefreshToken)
	}

	//auth route
	authGroup := apiGroup.Group("/auth")
	authGroup.Use(middleware.JwtAuth())
	{
		authGroup.GET("/user", authController.GetUserInfo)
		authGroup.POST("/change-password", authController.ChangePassword)
	}
}

// User routes
func setupUserRoutes(apiGroup *gin.RouterGroup) {
	userController := controller.NewUserController()

	//base routeGroup
	baseUserGroup := apiGroup.Group("/users")
	baseUserGroup.Use(middleware.JwtAuth())

	baseUserGroup.GET("", middleware.PermissionAuth("user:list"), userController.GetUsers)
	baseUserGroup.GET("/:id", middleware.AnyPermissionAuth("user:list", "user:read"), userController.GetUser)
	baseUserGroup.POST("", middleware.RoleOrPermissionAuth("admin", "user:create"), userController.CreateUser)
	baseUserGroup.PUT("/:id", middleware.RoleAndPermissionAuth("admin", "user:update"), userController.UpdateUser)
	baseUserGroup.DELETE("/:id", middleware.AllPermissionsAuth("user:delete", "user:manage"), userController.DeleteUser)

	//admin group
	adminGroup := apiGroup.Group("/admin/users")
	adminGroup.Use(middleware.JwtAuth())
	adminGroup.Use(middleware.RoleAuth("admin"))
	{ //delete user
		adminGroup.DELETE("/:id", userController.DeleteUser)
	}
}
