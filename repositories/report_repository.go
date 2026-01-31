package repositories

import (
	"time"

	entity "dinacom-11.0-backend/models/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepository interface {
	CreateReport(report *entity.Report) error
	GetCompletedNonGoodReports() ([]entity.Report, error)
	GetReportByID(id string) (*entity.Report, error)
	AssignWorker(reportID string, workerID uuid.UUID, adminNotes string, deadline *time.Time) error
	GetAssignedReports() ([]entity.Report, error)
	UpdateAfterImage(reportID string, afterImageURL string, status string) error
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

func (r *reportRepository) GetCompletedNonGoodReports() ([]entity.Report, error) {
	var reports []entity.Report
	err := r.db.Where("status = ? AND destruct_class != ?", "completed", "good").Find(&reports).Error
	return reports, err
}

func (r *reportRepository) GetReportByID(id string) (*entity.Report, error) {
	var report entity.Report
	err := r.db.Where("id = ?", id).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) AssignWorker(reportID string, workerID uuid.UUID, adminNotes string, deadline *time.Time) error {
	return r.db.Model(&entity.Report{}).Where("id = ?", reportID).Updates(map[string]interface{}{
		"worker_id":   workerID,
		"status":      "assigned",
		"admin_notes": adminNotes,
		"deadline":    deadline,
	}).Error
}

func (r *reportRepository) GetAssignedReports() ([]entity.Report, error) {
	var reports []entity.Report
	err := r.db.Where("worker_id IS NOT NULL").Find(&reports).Error
	return reports, err
}

func (r *reportRepository) UpdateAfterImage(reportID string, afterImageURL string, status string) error {
	return r.db.Model(&entity.Report{}).Where("id = ?", reportID).Updates(map[string]interface{}{
		"after_image_url": afterImageURL,
		"status":          status,
	}).Error
}
