package dto

import "time"

type UserReportResponse struct {
	ID             string     `json:"id"`
	Longitude      float64    `json:"longitude"`
	Latitude       float64    `json:"latitude"`
	RoadName       string     `json:"road_name"`
	BeforeImageURL string     `json:"before_image_url"`
	AfterImageURL  string     `json:"after_image_url"`
	Description    string     `json:"description"`
	DestructClass  string     `json:"destruct_class"`
	LocationScore  float64    `json:"location_score"`
	TotalScore     float64    `json:"total_score"`
	Status         string     `json:"status"`
	AdminNotes     string     `json:"admin_notes"`
	Deadline       *time.Time `json:"deadline"`
	CreatedAt      time.Time  `json:"created_at"`
}

type PaginatedReportsResponse struct {
	Reports    []UserReportResponse `json:"reports"`
	TotalCount int64                `json:"total_count"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

type VerifyReportRequest struct {
	ReportID string `json:"report_id" binding:"required"`
}
