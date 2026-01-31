package controllers

import (
	"net/http"

	"dinacom-11.0-backend/services"
	"dinacom-11.0-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RankController interface {
	GetUserRank(ctx *gin.Context)
}

type rankController struct {
	rankService services.RankService
}

func NewRankController(rankService services.RankService) RankController {
	return &rankController{rankService: rankService}
}

// @Summary Get User Rank
// @Description Get the logged-in user's rank, level, XP progress
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserRankResponse
// @Failure 401 {object} map[string]string
// @Router /api/user/rank [get]
func (c *rankController) GetUserRank(ctx *gin.Context) {
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

	response, err := c.rankService.GetUserRank(userID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "User rank retrieved", response)
}
