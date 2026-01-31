package router

import (
	_ "dinacom-11.0-backend/docs"
	"dinacom-11.0-backend/provider"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunRouter(appProvider provider.AppProvider) {
	router, controller, config := appProvider.ProvideRouter(), appProvider.ProvideControllers(), appProvider.ProvideConfig()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authRouter := NewAuthRouter(controller.ProvideAuthController())
	authRouter.Setup(router.Group("/api"))

	reportRouter := NewReportRouter(controller.ProvideReportController())
	reportRouter.Setup(router.Group("/api"))

	achievementRouter := NewAchievementRouter(controller.ProvideAchievementController(), controller.ProvideRankController())
	achievementRouter.Setup(router.Group("/api"))

	notificationRouter := NewNotificationRouter(controller.ProvideNotificationController())
	notificationRouter.Setup(router.Group("/api"))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(config.ProvideEnvConfig().GetTCPAddress())
	if err != nil {
		panic(err)
	}
}
