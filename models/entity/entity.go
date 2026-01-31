package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username  string         `gorm:"type:varchar(100);not null;unique" json:"username"`
	Fullname  string         `gorm:"column:name;type:varchar(100);not null" json:"fullname"`
	Email     string         `gorm:"type:varchar(100);not null;unique" json:"email"`
	Role      string         `gorm:"type:varchar(20);not null;default:'user'" json:"role"` // user, admin, worker
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Verified  bool           `gorm:"default:false" json:"verified"`
	TotalXP   int            `gorm:"default:0" json:"total_xp"`
	Level     int            `gorm:"default:1" json:"level"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Report struct {
	ID             string         `gorm:"type:text;primary_key" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid" json:"user_id"`
	WorkerID       *uuid.UUID     `gorm:"type:uuid" json:"worker_id"`
	Longitude      float64        `gorm:"type:numeric" json:"longitude"`
	Latitude       float64        `gorm:"type:numeric" json:"latitude"`
	RoadName       string         `gorm:"column:road_name;type:text" json:"road_name"`
	BeforeImageURL string         `gorm:"column:before_image_url;type:text" json:"before_image_url"`
	AfterImageURL  string         `gorm:"column:after_image_url;type:text" json:"after_image_url"`
	Description    string         `gorm:"type:text" json:"description"`
	DestructClass  string         `gorm:"column:destruct_class;type:text" json:"destruct_class"`
	LocationScore  float64        `gorm:"column:location_score;type:numeric" json:"location_score"`
	TotalScore     float64        `gorm:"column:total_score;type:numeric" json:"total_score"`
	Status         string         `gorm:"type:text" json:"status"`
	AdminNotes     string         `gorm:"column:admin_notes;type:text" json:"admin_notes"`
	Deadline       *time.Time     `gorm:"column:deadline;type:timestamp" json:"deadline"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Achievement struct {
	ID          string `gorm:"type:text;primary_key" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	BadgeURL    string `gorm:"column:badge_url;type:text" json:"badge_url"`
	Category    string `gorm:"type:varchar(50)" json:"category"` // milestone, quality, explorer, streak
	XPReward    int    `gorm:"column:xp_reward;default:0" json:"xp_reward"`
}

type UserAchievement struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	AchievementID string      `gorm:"type:text;not null" json:"achievement_id"`
	Achievement   Achievement `gorm:"foreignKey:AchievementID" json:"achievement"`
	UnlockedAt    time.Time   `json:"unlocked_at"`
}
