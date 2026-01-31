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
	router.GET("/get_report", r.reportController.GetReports)

	reportGroup := router.Group("/user/report")
	reportGroup.Use(middleware.AuthMiddleware())
	reportGroup.POST("", r.reportController.CreateReport)

	adminGroup := router.Group("/admin/report")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.RoleMiddleware("admin"))
	adminGroup.PATCH("/assign", r.reportController.AssignWorker)
	adminGroup.GET("/assign", r.reportController.GetAssignedReports)

	workerGroup := router.Group("/worker")
	workerGroup.Use(middleware.AuthMiddleware())
	workerGroup.Use(middleware.RoleMiddleware("worker", "admin"))
	workerGroup.PATCH("/report", r.reportController.FinishReport)
}
