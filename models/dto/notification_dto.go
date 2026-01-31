package dto

import "github.com/google/uuid"

// NotificationResponse represents a single notification for frontend
type NotificationResponse struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Data      string    `json:"data,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt string    `json:"created_at"`
}

// NotificationListResponse wraps list of notifications with unread count
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                  `json:"unread_count"`
	Total         int64                  `json:"total"`
}

// MarkNotificationReadRequest for marking notification as read
type MarkNotificationReadRequest struct {
	NotificationID uuid.UUID `json:"notification_id" binding:"required"`
}

// CreateNotificationRequest for admin/system to create notifications
type CreateNotificationRequest struct {
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Type    string    `json:"type" binding:"required"`
	Title   string    `json:"title" binding:"required"`
	Message string    `json:"message" binding:"required"`
	Data    string    `json:"data"`
}
