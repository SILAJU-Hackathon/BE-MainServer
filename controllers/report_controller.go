package controllers

import (
	"encoding/json"
	"net/http"

	"dinacom-11.0-backend/models/dto"
	"dinacom-11.0-backend/services"
	"dinacom-11.0-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportController interface {
	CreateReport(ctx *gin.Context)
}

type reportController struct {
	reportService services.ReportService
}

func NewReportController(reportService services.ReportService) ReportController {
	return &reportController{reportService: reportService}
}

// @Summary Create Report
// @Description Submit a new report with image and location data
// @Tags Report
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Image file (JPG, PNG, JPEG, max 32MB)"
// @Param json formData string true "JSON data: {\"longitude\": 0, \"latitude\": 0, \"road_name\": \"string\", \"description\": \"string\"}"
// @Security BearerAuth
// @Success 200 {object} dto.ReportResponse
// @Failure 400 {object} map[string]string
// @Router /api/report [post]
func (c *reportController) CreateReport(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	file, header, err := ctx.Request.FormFile("files")
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Image file is required")
		return
	}
	defer file.Close()

	jsonData := ctx.PostForm("json")
	if jsonData == "" {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "JSON data is required")
		return
	}

	var req dto.ReportRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	response, err := c.reportService.CreateReport(userID, file, header, req)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Report created successfully", response)
}
