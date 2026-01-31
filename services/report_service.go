package services

import (
	"errors"
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
	GetReports() ([]dto.ReportLocationResponse, error)
	AssignWorker(req dto.AssignWorkerRequest) (string, error)
	GetAssignedReports() ([]dto.AssignedWorkerResponse, error)
	FinishReport(workerID uuid.UUID, file multipart.File, header *multipart.FileHeader, reportID string) error
}

type reportService struct {
	reportRepo       repositories.ReportRepository
	userRepo         repositories.UserRepository
	cloudinaryClient *utils.CloudinaryClient
}

func NewReportService(reportRepo repositories.ReportRepository, userRepo repositories.UserRepository) ReportService {
	client, _ := utils.NewCloudinaryClient()
	return &reportService{
		reportRepo:       reportRepo,
		userRepo:         userRepo,
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

func (s *reportService) GetReports() ([]dto.ReportLocationResponse, error) {
	reports, err := s.reportRepo.GetCompletedNonGoodReports()
	if err != nil {
		return nil, err
	}

	var response []dto.ReportLocationResponse
	for _, report := range reports {
		response = append(response, dto.ReportLocationResponse{
			ID:            report.ID,
			Longitude:     report.Longitude,
			Latitude:      report.Latitude,
			DestructClass: report.DestructClass,
			TotalScore:    report.TotalScore,
		})
	}

	return response, nil
}

func (s *reportService) AssignWorker(req dto.AssignWorkerRequest) (string, error) {
	report, err := s.reportRepo.GetReportByID(req.ReportID)
	if err != nil {
		return "", errors.New("report not found")
	}

	if report.WorkerID != nil {
		return "", errors.New("report already assigned to a worker")
	}

	worker, err := s.userRepo.FindUserByID(req.WorkerID)
	if err != nil || worker == nil {
		return "", errors.New("worker not found")
	}

	if worker.Role != "worker" {
		return "", errors.New("only workers can be assigned")
	}

	if err := s.reportRepo.AssignWorker(req.ReportID, req.WorkerID, req.AdminNotes, req.Deadline); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s is successfully assigned", worker.Fullname), nil
}

func (s *reportService) GetAssignedReports() ([]dto.AssignedWorkerResponse, error) {
	reports, err := s.reportRepo.GetAssignedReports()
	if err != nil {
		return nil, err
	}

	var response []dto.AssignedWorkerResponse
	for _, report := range reports {
		workerName := ""
		if report.WorkerID != nil {
			worker, _ := s.userRepo.FindUserByID(*report.WorkerID)
			if worker != nil {
				workerName = worker.Fullname
			}
		}

		response = append(response, dto.AssignedWorkerResponse{
			WorkerName: workerName,
			RoadName:   report.RoadName,
			Longitude:  report.Longitude,
			Latitude:   report.Latitude,
			Status:     report.Status,
		})
	}

	return response, nil
}

func (s *reportService) FinishReport(workerID uuid.UUID, file multipart.File, header *multipart.FileHeader, reportID string) error {
	if header.Size > maxFileSize {
		return http_error.FILE_TOO_LARGE
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return http_error.INVALID_FILE_FORMAT
	}

	report, err := s.reportRepo.GetReportByID(reportID)
	if err != nil {
		return errors.New("report not found")
	}

	if report.WorkerID == nil || *report.WorkerID != workerID {
		return errors.New("you are not assigned to this report")
	}

	afterImageID := fmt.Sprintf("%s_after_%s", reportID, time.Now().Format("20060102150405"))
	afterImageURL, err := s.cloudinaryClient.UploadImage(file, afterImageID, report.Longitude, report.Latitude, "After image")
	if err != nil {
		return http_error.CLOUDINARY_UPLOAD_FAILED
	}

	return s.reportRepo.UpdateAfterImage(reportID, afterImageURL, "Finish by Worker")
}
