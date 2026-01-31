package controllers

import (
	"net/http"

	"dinacom-11.0-backend/services"
	"dinacom-11.0-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AchievementController interface {
	GetUserAchievements(ctx *gin.Context)
	GetUnlockedAchievements(ctx *gin.Context)
	CheckAchievements(ctx *gin.Context)
}

type achievementController struct {
	achievementService services.AchievementService
}

func NewAchievementController(achievementService services.AchievementService) AchievementController {
	return &achievementController{achievementService: achievementService}
}

// @Summary Get All Achievements
// @Description Get all achievements with unlock status for the logged-in user
// @Tags Achievements
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AchievementListResponse
// @Failure 401 {object} map[string]string
// @Router /api/user/achievements [get]
func (c *achievementController) GetUserAchievements(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	response, err := c.achievementService.GetUserAchievements(userID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Achievements retrieved", response)
}

// @Summary Get Unlocked Achievements
// @Description Get only unlocked achievements for the logged-in user
// @Tags Achievements
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AchievementListResponse
// @Failure 401 {object} map[string]string
// @Router /api/user/achievements/unlocked [get]
func (c *achievementController) GetUnlockedAchievements(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	response, err := c.achievementService.GetUnlockedAchievements(userID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Unlocked achievements retrieved", response)
}

// @Summary Check Achievements
// @Description Trigger achievement check and unlock any new achievements
// @Tags Achievements
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []dto.NewAchievementResponse
// @Failure 401 {object} map[string]string
// @Router /api/user/achievements/check [post]
func (c *achievementController) CheckAchievements(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	newAchievements, err := c.achievementService.CheckAndUnlockAchievements(userID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if len(newAchievements) == 0 {
		utils.SendSuccessResponse(ctx, "No new achievements unlocked", []interface{}{})
		return
	}

	utils.SendSuccessResponse(ctx, "New achievements unlocked!", newAchievements)
}
