package router

import (
	"dinacom-11.0-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/gzip"
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
	{
		userGroup := authGroup.Group("/user")
		{
			userGroup.POST("/register", r.authController.RegisterUser)
			userGroup.POST("/verify-otp", r.authController.VerifyOTP)
			userGroup.POST("/login", r.authController.LoginUser)
		}

		adminGroup := authGroup.Group("/admin")
		{
			adminGroup.POST("/login", r.authController.LoginAdmin)
		}

		workerGroup := authGroup.Group("/worker")
		{
			workerGroup.POST("/login", r.authController.LoginWorker)
		}
	}
}
