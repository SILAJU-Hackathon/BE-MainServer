package router

import (
	"dinacom-11.0-backend/controllers"
	"dinacom-11.0-backend/middleware"
	"dinacom-11.0-backend/models/entity"

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

	userReportGroup := router.Group("/user/report")
	userReportGroup.Use(middleware.AuthMiddleware())
	userReportGroup.POST("", r.reportController.CreateReport)
	userReportGroup.GET("/me", r.reportController.GetUserReports)
	userReportGroup.GET("/stats", r.reportController.GetUserReportStats)

	adminGroup := router.Group("/admin/report")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.RoleMiddleware(entity.ROLE_ADMIN))
	adminGroup.GET("", r.reportController.GetAllReportsAdmin)
	adminGroup.GET("/detail", r.reportController.GetFullReportDetail)
	adminGroup.PATCH("/assign", r.reportController.AssignWorker)
	adminGroup.GET("/assign", r.reportController.GetAssignedReports)
	adminGroup.PATCH("/verify", r.reportController.VerifyReport)

	workerGroup := router.Group("/worker")
	workerGroup.Use(middleware.AuthMiddleware())
	workerGroup.Use(middleware.RoleMiddleware(entity.ROLE_WORKER, entity.ROLE_ADMIN))
	workerGroup.PATCH("/report", r.reportController.FinishReport)
	workerGroup.GET("/report/assign/me", r.reportController.GetWorkerAssignedReports)
	workerGroup.GET("/report/assign/detail", r.reportController.GetReportDetail)
	workerGroup.GET("/report/assign/image", r.reportController.GetReportImage)
	workerGroup.GET("/report/history/me", r.reportController.GetWorkerHistory)
}
