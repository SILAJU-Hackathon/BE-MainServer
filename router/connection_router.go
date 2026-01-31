package router

import (
	"dinacom-11.0-backend/provider"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func ConnectionRouter(router *gin.Engine, controller provider.ControllerProvider) {
	connectionController := controller.ProvideConnectionController()
	routerGroup := router.Group("/api")
	routerGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		routerGroup.POST("/connect", connectionController.Connect)
	}

}
