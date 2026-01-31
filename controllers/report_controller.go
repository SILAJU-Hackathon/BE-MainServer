package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	GetUserReports(ctx *gin.Context)
	GetWorkerAssignedReports(ctx *gin.Context)
	GetWorkerHistory(ctx *gin.Context)
	VerifyReport(ctx *gin.Context)
	GetReportDetail(ctx *gin.Context)
	GetReportImage(ctx *gin.Context)
	GetAllReportsAdmin(ctx *gin.Context)
	GetFullReportDetail(ctx *gin.Context)
	GetUserReportStats(ctx *gin.Context)
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
// @Param json formData string true "JSON data" default({"longitude": 106.816666, "latitude": -6.200000, "road_name": "Jalan Sudirman", "description": "Lubang besar di tengah jalan"})
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
// @Param json formData string true "JSON data" default({"report_id": "uuid-here"})
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

// @Summary Get User's Reports
// @Description Get all reports created by the logged-in user with pagination
// @Tags User
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} dto.PaginatedReportsResponse
// @Failure 401 {object} map[string]string
// @Router /api/user/report/me [get]
func (c *reportController) GetUserReports(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	response, err := c.reportService.GetUserReports(userID, page, limit)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "User reports retrieved", response)
}

// @Summary Get Worker's Assigned Reports
// @Description Get all reports assigned to the logged-in worker with pagination
// @Tags Worker
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} dto.PaginatedReportsResponse
// @Failure 401 {object} map[string]string
// @Router /api/worker/report/assign/me [get]
func (c *reportController) GetWorkerAssignedReports(ctx *gin.Context) {
	workerIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	workerID := workerIDVal.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	response, err := c.reportService.GetWorkerAssignedReports(workerID, page, limit)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Worker assigned reports retrieved", response)
}

// @Summary Get Worker's History
// @Description Get worker's completed reports with pagination and status filter
// @Tags Worker
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param verify_admin query bool false "True: finished, False: Finish by Worker" default(false)
// @Security BearerAuth
// @Success 200 {object} dto.PaginatedReportsResponse
// @Failure 401 {object} map[string]string
// @Router /api/worker/report/history/me [get]
func (c *reportController) GetWorkerHistory(ctx *gin.Context) {
	workerIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	workerID := workerIDVal.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	verifyAdmin := ctx.DefaultQuery("verify_admin", "false") == "true"
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	response, err := c.reportService.GetWorkerHistory(workerID, verifyAdmin, page, limit)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Worker history retrieved", response)
}

// @Summary Verify Report by Admin
// @Description Admin verifies a report with status 'Finish by Worker' to 'finished'
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body dto.VerifyReportRequest true "Verify Report Request"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/admin/report/verify [patch]
func (c *reportController) VerifyReport(ctx *gin.Context) {
	var req dto.VerifyReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.reportService.VerifyReport(req.ReportID); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Report verified successfully", nil)
}

// @Summary Get Report Detail
// @Description Get assigned report detail for worker
// @Tags Worker
// @Produce json
// @Param report_id query string true "Report ID"
// @Security BearerAuth
// @Success 200 {object} dto.ReportDetailResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/worker/report/assign/detail [get]
func (c *reportController) GetReportDetail(ctx *gin.Context) {
	workerIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	workerID := workerIDVal.(uuid.UUID)

	reportID := ctx.Query("report_id")
	if reportID == "" {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "report_id is required")
		return
	}

	response, err := c.reportService.GetReportDetail(workerID, reportID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Report detail retrieved", response)
}

// @Summary Get Report Image
// @Description Get before image URL for assigned report
// @Tags Worker
// @Produce json
// @Param report_id query string true "Report ID"
// @Security BearerAuth
// @Success 200 {object} dto.ReportImageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/worker/report/assign/image [get]
func (c *reportController) GetReportImage(ctx *gin.Context) {
	workerIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	workerID := workerIDVal.(uuid.UUID)

	reportID := ctx.Query("report_id")
	if reportID == "" {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "report_id is required")
		return
	}

	response, err := c.reportService.GetReportImage(workerID, reportID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Report image retrieved", response)
}

// @Summary Get All Reports (Admin)
// @Description Get all reports with optional status filter and pagination
// @Tags Admin
// @Produce json
// @Param status query string false "Filter by status (pending, assigned, finish by worker, verified, complete)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} dto.PaginatedReportsResponse
// @Failure 401 {object} map[string]string
// @Router /api/admin/report [get]
func (c *reportController) GetAllReportsAdmin(ctx *gin.Context) {
	status := ctx.Query("status")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	response, err := c.reportService.GetAllReportsAdmin(status, page, limit)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Reports retrieved", response)
}

// @Summary Get Full Report Detail (Admin)
// @Description Get complete report details by report ID
// @Tags Admin
// @Produce json
// @Param report_id query string true "Report ID"
// @Security BearerAuth
// @Success 200 {object} dto.FullReportDetailResponse
// @Failure 400 {object} map[string]string
// @Router /api/admin/report/detail [get]
func (c *reportController) GetFullReportDetail(ctx *gin.Context) {
	reportID := ctx.Query("report_id")
	if reportID == "" {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "report_id is required")
		return
	}

	response, err := c.reportService.GetFullReportDetail(reportID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Report detail retrieved", response)
}

// @Summary Get User Report Stats
// @Description Get report statistics for the authenticated user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserReportStatsResponse
// @Failure 401 {object} map[string]string
// @Router /api/user/report/stats [get]
func (c *reportController) GetUserReportStats(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID := userIDVal.(uuid.UUID)
	stats, err := c.reportService.GetUserReportStats(userID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Report stats retrieved", stats)
}
