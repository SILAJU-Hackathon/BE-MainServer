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
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
