package repositories

import (
	"time"

	"dinacom-11.0-backend/models/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AchievementRepository interface {
	GetAllAchievements() ([]entity.Achievement, error)
	GetUserAchievements(userID uuid.UUID) ([]entity.UserAchievement, error)
	HasAchievement(userID uuid.UUID, achievementID string) (bool, error)
	UnlockAchievement(userID uuid.UUID, achievementID string) error
	GetUserReportCount(userID uuid.UUID) (int64, error)
	GetUserVerifiedReportCount(userID uuid.UUID) (int64, error)
	GetUserUniqueRoadCount(userID uuid.UUID) (int64, error)
	HasHighImpactReport(userID uuid.UUID) (bool, error)
	HasCriticalReport(userID uuid.UUID) (bool, error)
	GetUserReportDatesLast7Days(userID uuid.UUID) ([]time.Time, error)
	HasEarlyMorningReport(userID uuid.UUID) (bool, error)
}

type achievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) GetAllAchievements() ([]entity.Achievement, error) {
	var achievements []entity.Achievement
	err := r.db.Find(&achievements).Error
	return achievements, err
}

func (r *achievementRepository) GetUserAchievements(userID uuid.UUID) ([]entity.UserAchievement, error) {
	var userAchievements []entity.UserAchievement
	err := r.db.Preload("Achievement").Where("user_id = ?", userID).Find(&userAchievements).Error
	return userAchievements, err
}

func (r *achievementRepository) HasAchievement(userID uuid.UUID, achievementID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.UserAchievement{}).Where("user_id = ? AND achievement_id = ?", userID, achievementID).Count(&count).Error
	return count > 0, err
}

func (r *achievementRepository) UnlockAchievement(userID uuid.UUID, achievementID string) error {
	userAchievement := entity.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
		UnlockedAt:    time.Now(),
	}
	return r.db.Create(&userAchievement).Error
}

func (r *achievementRepository) GetUserReportCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Report{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *achievementRepository) GetUserVerifiedReportCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Report{}).Where("user_id = ? AND status = ?", userID, entity.STATUS_VERIFIED).Count(&count).Error
	return count, err
}

func (r *achievementRepository) GetUserUniqueRoadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Report{}).Where("user_id = ?", userID).Distinct("road_name").Count(&count).Error
	return count, err
}

func (r *achievementRepository) HasHighImpactReport(userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Report{}).Where("user_id = ? AND total_score > ?", userID, 80).Count(&count).Error
	return count > 0, err
}

func (r *achievementRepository) HasCriticalReport(userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Report{}).Where("user_id = ? AND destruct_class = ?", userID, "Very Poor").Count(&count).Error
	return count > 0, err
}

func (r *achievementRepository) GetUserReportDatesLast7Days(userID uuid.UUID) ([]time.Time, error) {
	var dates []time.Time
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	err := r.db.Model(&entity.Report{}).
		Where("user_id = ? AND created_at >= ?", userID, sevenDaysAgo).
		Pluck("DATE(created_at)", &dates).Error
	return dates, err
}

func (r *achievementRepository) HasEarlyMorningReport(userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Report{}).
		Where("user_id = ? AND EXTRACT(HOUR FROM created_at) < 7", userID).
		Count(&count).Error
	return count > 0, err
}
