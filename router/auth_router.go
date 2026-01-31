package router

import (
	"dinacom-11.0-backend/controllers"
	"dinacom-11.0-backend/middleware"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type AuthRouter interface {
	Setup(router *gin.RouterGroup)
}

type authRouter struct {
	authController controllers.AuthController
}

func NewAuthRouter(authController controllers.AuthController) AuthRouter {
	return &authRouter{authController: authController}
}

func (r *authRouter) Setup(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	authGroup.Use(gzip.Gzip(gzip.DefaultCompression))

	userGroup := authGroup.Group("/user")
	userGroup.POST("/register", r.authController.RegisterUser)
	userGroup.POST("/verify-otp", r.authController.VerifyOTP)
	userGroup.POST("/login", r.authController.LoginUser)

	adminGroup := authGroup.Group("/admin")
	adminGroup.POST("/login", r.authController.LoginAdmin)

	workerGroup := authGroup.Group("/worker")
	workerGroup.POST("/login", r.authController.LoginWorker)

	protectedGroup := authGroup.Group("")
	protectedGroup.Use(middleware.AuthMiddleware())
	protectedGroup.GET("/me", r.authController.GetProfile)

	adminProtected := authGroup.Group("/admin")
	adminProtected.Use(middleware.AuthMiddleware())
	adminProtected.Use(middleware.RoleMiddleware("admin"))
	adminProtected.GET("/users", r.authController.GetAllUsers)
	adminProtected.GET("/workers", r.authController.GetAllWorkers)

	workerProtected := authGroup.Group("/worker")
	workerProtected.Use(middleware.AuthMiddleware())
	workerProtected.Use(middleware.RoleMiddleware("worker", "admin"))
	workerProtected.GET("/me", r.authController.GetProfile)

	userProtected := authGroup.Group("/user")
	userProtected.Use(middleware.AuthMiddleware())
	userProtected.Use(middleware.RoleMiddleware("user", "admin"))
	userProtected.GET("/me", r.authController.GetProfile)
}
