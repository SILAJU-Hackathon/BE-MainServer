package repositories

import (
	entity "dinacom-11.0-backend/models/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *entity.Notification) error
	GetByUserID(userID uuid.UUID, limit, offset int) ([]entity.Notification, int64, error)
	GetUnreadCount(userID uuid.UUID) (int64, error)
	MarkAsRead(notificationID uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	Delete(notificationID uuid.UUID) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]entity.Notification, int64, error) {
	var notifications []entity.Notification
	var total int64

	r.db.Model(&entity.Notification{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *notificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).
		Where("id = ?", notificationID).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (r *notificationRepository) Delete(notificationID uuid.UUID) error {
	return r.db.Delete(&entity.Notification{}, "id = ?", notificationID).Error
}
