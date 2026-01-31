package router

import (
	"dinacom-11.0-backend/controllers"
	"dinacom-11.0-backend/middleware"

	"github.com/gin-gonic/gin"
)

type NotificationRouter interface {
	Setup(router *gin.RouterGroup)
}

type notificationRouter struct {
	notifController controllers.NotificationController
}

func NewNotificationRouter(notifController controllers.NotificationController) NotificationRouter {
	return &notificationRouter{notifController: notifController}
}

func (r *notificationRouter) Setup(router *gin.RouterGroup) {
	notifGroup := router.Group("/user/notifications")
	notifGroup.Use(middleware.AuthMiddleware())

	notifGroup.GET("", r.notifController.GetNotifications)
	notifGroup.GET("/unread-count", r.notifController.GetUnreadCount)
	notifGroup.PATCH("/read", r.notifController.MarkAsRead)
	notifGroup.PATCH("/read-all", r.notifController.MarkAllAsRead)
}
