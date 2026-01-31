package router

import (
	"dinacom-11.0-backend/controllers"
	"dinacom-11.0-backend/middleware"

	"github.com/gin-gonic/gin"
)

type ReportRouter interface {
	Setup(router *gin.RouterGroup)
}

type reportRouter struct {
	reportController controllers.ReportController
}

func NewReportRouter(reportController controllers.ReportController) ReportRouter {
	return &reportRouter{reportController: reportController}
}

func (r *reportRouter) Setup(router *gin.RouterGroup) {
	reportGroup := router.Group("/user/report")
	reportGroup.Use(middleware.AuthMiddleware())
	reportGroup.POST("", r.reportController.CreateReport)
}
