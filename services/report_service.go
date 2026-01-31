package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"dinacom-11.0-backend/models/dto"
	entity "dinacom-11.0-backend/models/entity"
	http_error "dinacom-11.0-backend/models/error"
	"dinacom-11.0-backend/repositories"
	"dinacom-11.0-backend/utils"

	"github.com/google/uuid"
)

const maxFileSize = 32 << 20

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

type ReportService interface {
	CreateReport(userID uuid.UUID, file multipart.File, header *multipart.FileHeader, req dto.ReportRequest) (*dto.ReportResponse, error)
}

type reportService struct {
	reportRepo       repositories.ReportRepository
	cloudinaryClient *utils.CloudinaryClient
}

func NewReportService(reportRepo repositories.ReportRepository) ReportService {
	client, _ := utils.NewCloudinaryClient()
	return &reportService{
		reportRepo:       reportRepo,
		cloudinaryClient: client,
	}
}

func (s *reportService) CreateReport(userID uuid.UUID, file multipart.File, header *multipart.FileHeader, req dto.ReportRequest) (*dto.ReportResponse, error) {
	if header.Size > maxFileSize {
		return nil, http_error.FILE_TOO_LARGE
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return nil, http_error.INVALID_FILE_FORMAT
	}

	reportID := fmt.Sprintf("%s_%s_%.6f_%.6f", uuid.New().String(), time.Now().Format("20060102150405"), req.Longitude, req.Latitude)

	imageURL, err := s.cloudinaryClient.UploadImage(file, reportID, req.Longitude, req.Latitude, req.Description)
	if err != nil {
		return nil, http_error.CLOUDINARY_UPLOAD_FAILED
	}

	report := &entity.Report{
		ID:             reportID,
		UserID:         userID,
		Longitude:      req.Longitude,
		Latitude:       req.Latitude,
		RoadName:       req.RoadName,
		BeforeImageURL: imageURL,
		Description:    req.Description,
		Status:         "pending",
	}

	if err := s.reportRepo.CreateReport(report); err != nil {
		return nil, http_error.REPORT_CREATION_FAILED
	}

	return &dto.ReportResponse{
		ID:             report.ID,
		UserID:         report.UserID,
		Longitude:      report.Longitude,
		Latitude:       report.Latitude,
		RoadName:       report.RoadName,
		BeforeImageURL: report.BeforeImageURL,
		AfterImageURL:  report.AfterImageURL,
		Description:    report.Description,
		DestructClass:  report.DestructClass,
		LocationScore:  report.LocationScore,
		TotalScore:     report.TotalScore,
		Status:         report.Status,
	}, nil
}
