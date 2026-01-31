package dto

type ReportLocationResponse struct {
	ID            string  `json:"id"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	DestructClass string  `json:"destruct_class"`
	TotalScore    float64 `json:"total_score"`
}
