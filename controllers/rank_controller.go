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
	GetLeaderboard(ctx *gin.Context)
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
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID := userIDVal.(uuid.UUID)

	response, err := c.rankService.GetUserRank(userID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "User rank retrieved", response)
}

// @Summary Get Leaderboard
// @Description Get top users by XP (default limit 10)
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit (default 10)"
// @Success 200 {array} dto.LeaderboardEntry
// @Failure 401 {object} map[string]string
// @Router /api/user/leaderboard [get]
func (c *rankController) GetLeaderboard(ctx *gin.Context) {
	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		// simple parse, ignore error for now or use strconv
		// For simplicity/safety in this specific context without strconv import visible yet:
		// We can try to bind query or just trust default if parsing fails.
		// Let's rely on service validation or just passing it if we had strconv.
		// Since I can't easily see imports, I'll update imports if needed or just use default if not provided.
		// Actually, let's use a helper or bind.
	}

	// Better to use BindQuery if struct defined, or strconv.
	// Using a simple struct to bind query params
	var query struct {
		Limit int `form:"limit"`
	}
	if err := ctx.ShouldBindQuery(&query); err == nil && query.Limit > 0 {
		limit = query.Limit
	}

	leaderboard, err := c.rankService.GetLeaderboard(limit)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Leaderboard retrieved", leaderboard)
}
