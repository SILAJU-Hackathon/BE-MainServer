package repositories

import (
	entity "dinacom-11.0-backend/models/entity"

	"gorm.io/gorm"
)

type ReportRepository interface {
	CreateReport(report *entity.Report) error
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) CreateReport(report *entity.Report) error {
	return r.db.Create(report).Error
}
