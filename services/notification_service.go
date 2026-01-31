package services

import (
	"dinacom-11.0-backend/models/dto"
	entity "dinacom-11.0-backend/models/entity"
	"dinacom-11.0-backend/repositories"

	"github.com/google/uuid"
)

type NotificationService interface {
	GetUserNotifications(userID uuid.UUID, page, limit int) (*dto.NotificationListResponse, error)
	MarkAsRead(userID, notificationID uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	CreateNotification(userID uuid.UUID, notifType, title, message, data string) error
	GetUnreadCount(userID uuid.UUID) (int64, error)
}

type notificationService struct {
	notifRepo repositories.NotificationRepository
}

func NewNotificationService(notifRepo repositories.NotificationRepository) NotificationService {
	return &notificationService{notifRepo: notifRepo}
}

func (s *notificationService) GetUserNotifications(userID uuid.UUID, page, limit int) (*dto.NotificationListResponse, error) {
	offset := (page - 1) * limit

	notifications, total, err := s.notifRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.notifRepo.GetUnreadCount(userID)
	if err != nil {
		return nil, err
	}

	var response []dto.NotificationResponse
	for _, n := range notifications {
		response = append(response, dto.NotificationResponse{
			ID:        n.ID,
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			Data:      n.Data,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.NotificationListResponse{
		Notifications: response,
		UnreadCount:   unreadCount,
		Total:         total,
	}, nil
}

func (s *notificationService) MarkAsRead(userID, notificationID uuid.UUID) error {
	return s.notifRepo.MarkAsRead(notificationID)
}

func (s *notificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notifRepo.MarkAllAsRead(userID)
}

func (s *notificationService) CreateNotification(userID uuid.UUID, notifType, title, message, data string) error {
	notification := &entity.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		Data:    data,
		IsRead:  false,
	}
	return s.notifRepo.Create(notification)
}

func (s *notificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.notifRepo.GetUnreadCount(userID)
}
