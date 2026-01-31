package controllers

import (
	"net/http"
	"strconv"

	"dinacom-11.0-backend/models/dto"
	"dinacom-11.0-backend/services"
	"dinacom-11.0-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationController interface {
	GetNotifications(ctx *gin.Context)
	MarkAsRead(ctx *gin.Context)
	MarkAllAsRead(ctx *gin.Context)
	GetUnreadCount(ctx *gin.Context)
}

type notificationController struct {
	notifService services.NotificationService
}

func NewNotificationController(notifService services.NotificationService) NotificationController {
	return &notificationController{notifService: notifService}
}

// @Summary Get User Notifications
// @Description Get paginated notifications for authenticated user
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} dto.NotificationListResponse
// @Failure 401 {object} map[string]string
// @Router /api/user/notifications [get]
func (c *notificationController) GetNotifications(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	response, err := c.notifService.GetUserNotifications(userID, page, limit)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Notifications retrieved", response)
}

// @Summary Mark Notification as Read
// @Description Mark a specific notification as read
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.MarkNotificationReadRequest true "Notification ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/user/notifications/read [patch]
func (c *notificationController) MarkAsRead(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req dto.MarkNotificationReadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.notifService.MarkAsRead(userID, req.NotificationID); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Notification marked as read", nil)
}

// @Summary Mark All Notifications as Read
// @Description Mark all notifications as read for authenticated user
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/user/notifications/read-all [patch]
func (c *notificationController) MarkAllAsRead(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	if err := c.notifService.MarkAllAsRead(userID); err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "All notifications marked as read", nil)
}

// @Summary Get Unread Count
// @Description Get count of unread notifications
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]int64
// @Failure 401 {object} map[string]string
// @Router /api/user/notifications/unread-count [get]
func (c *notificationController) GetUnreadCount(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	count, err := c.notifService.GetUnreadCount(userID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Unread count retrieved", map[string]int64{"unread_count": count})
}
