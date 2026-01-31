package dto

import (
	"time"

	"github.com/google/uuid"
)

type AssignWorkerRequest struct {
	ReportID   string     `json:"report_id" binding:"required"`
	WorkerID   uuid.UUID  `json:"worker_id" binding:"required"`
	AdminNotes string     `json:"admin_notes"`
	Deadline   *time.Time `json:"deadline"`
}

type AssignedWorkerResponse struct {
	WorkerName string  `json:"worker_name"`
	RoadName   string  `json:"road_name"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Status     string  `json:"status"`
}

type WorkerReportRequest struct {
	ReportID string `json:"report_id" binding:"required"`
}
