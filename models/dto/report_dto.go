package dto

import "github.com/google/uuid"

type ReportRequest struct {
	Longitude   float64 `json:"longitude" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	RoadName    string  `json:"road_name" binding:"required"`
	Description string  `json:"description"`
}

type ReportResponse struct {
	ID             string    `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Longitude      float64   `json:"longitude"`
	Latitude       float64   `json:"latitude"`
	RoadName       string    `json:"road_name"`
	BeforeImageURL string    `json:"before_image_url"`
	AfterImageURL  string    `json:"after_image_url"`
	Description    string    `json:"description"`
	DestructClass  string    `json:"destruct_class"`
	LocationScore  float64   `json:"location_score"`
	TotalScore     float64   `json:"total_score"`
	Status         string    `json:"status"`
}
