package router

import (
	_ "dinacom-11.0-backend/docs"
	"dinacom-11.0-backend/provider"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunRouter(appProvider provider.AppProvider) {
	router, controller, config := appProvider.ProvideRouter(), appProvider.ProvideControllers(), appProvider.ProvideConfig()

	authRouter := NewAuthRouter(controller.ProvideAuthController())
	authRouter.Setup(router.Group("/api"))

	reportRouter := NewReportRouter(controller.ProvideReportController())
	reportRouter.Setup(router.Group("/api"))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(config.ProvideEnvConfig().GetTCPAddress())
	if err != nil {
		panic(err)
	}
}
