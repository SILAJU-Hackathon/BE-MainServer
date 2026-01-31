package router

import (
	"dinacom-11.0-backend/controllers"
	"dinacom-11.0-backend/middleware"

	"github.com/gin-gonic/gin"
)

type AchievementRouter interface {
	Setup(router *gin.RouterGroup)
}

type achievementRouter struct {
	achievementController controllers.AchievementController
}

func NewAchievementRouter(achievementController controllers.AchievementController) AchievementRouter {
	return &achievementRouter{achievementController: achievementController}
}

func (r *achievementRouter) Setup(router *gin.RouterGroup) {
	achievementGroup := router.Group("/user/achievements")
	achievementGroup.Use(middleware.AuthMiddleware())
	achievementGroup.GET("", r.achievementController.GetUserAchievements)
	achievementGroup.GET("/unlocked", r.achievementController.GetUnlockedAchievements)
	achievementGroup.POST("/check", r.achievementController.CheckAchievements)
}
