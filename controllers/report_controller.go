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
	GetReports(ctx *gin.Context)
	AssignWorker(ctx *gin.Context)
	GetAssignedReports(ctx *gin.Context)
	FinishReport(ctx *gin.Context)
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
// @Router /api/user/report [post]
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

// @Summary Get Reports
// @Description Get all completed reports with non-good destruct class
// @Tags Report
// @Produce json
// @Success 200 {array} dto.ReportLocationResponse
// @Failure 500 {object} map[string]string
// @Router /api/get_report [get]
func (c *reportController) GetReports(ctx *gin.Context) {
	reports, err := c.reportService.GetReports()
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Reports retrieved", reports)
}

// @Summary Assign Worker to Report
// @Description Admin assigns a worker to a report
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body dto.AssignWorkerRequest true "Assign Worker Request"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/admin/report/assign [patch]
func (c *reportController) AssignWorker(ctx *gin.Context) {
	var req dto.AssignWorkerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	message, err := c.reportService.AssignWorker(req)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, message, nil)
}

// @Summary Get Assigned Workers
// @Description Get all workers with assigned reports
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AssignedWorkerResponse
// @Failure 500 {object} map[string]string
// @Router /api/admin/report/assign [get]
func (c *reportController) GetAssignedReports(ctx *gin.Context) {
	reports, err := c.reportService.GetAssignedReports()
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Assigned reports retrieved", reports)
}

// @Summary Finish Report by Worker
// @Description Worker uploads after image and marks report as finished
// @Tags Worker
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "After image file"
// @Param json formData string true "JSON data: {\"report_id\": \"string\"}"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/worker/report [patch]
func (c *reportController) FinishReport(ctx *gin.Context) {
	workerIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	workerID := workerIDVal.(uuid.UUID)

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

	var req dto.WorkerReportRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := c.reportService.FinishReport(workerID, file, header, req.ReportID); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Report finished successfully", nil)
}
