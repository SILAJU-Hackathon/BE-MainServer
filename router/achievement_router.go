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
	rankController        controllers.RankController
}

func NewAchievementRouter(achievementController controllers.AchievementController, rankController controllers.RankController) AchievementRouter {
	return &achievementRouter{
		achievementController: achievementController,
		rankController:        rankController,
	}
}

func (r *achievementRouter) Setup(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())

	// Achievement routes
	achievementGroup := userGroup.Group("/achievements")
	achievementGroup.GET("", r.achievementController.GetUserAchievements)
	achievementGroup.GET("/unlocked", r.achievementController.GetUnlockedAchievements)
	achievementGroup.POST("/check", r.achievementController.CheckAchievements)

	// Rank route
	userGroup.GET("/rank", r.rankController.GetUserRank)
	userGroup.GET("/leaderboard", r.rankController.GetLeaderboard)
}
