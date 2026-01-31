package controllers

import (
	"net/http"

	"dinacom-11.0-backend/models/dto"
	"dinacom-11.0-backend/services"
	"dinacom-11.0-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController interface {
	RegisterUser(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	LoginAdmin(ctx *gin.Context)
	LoginWorker(ctx *gin.Context)
	GoogleAuth(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetAllWorkers(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService: authService}
}

// @Summary Register a new user
// @Description Register a new user with email verification via OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/auth/user/register [post]
func (c *authController) RegisterUser(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.authService.RegisterUser(req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "OTP sent to email", nil)
}

// @Summary Verify OTP
// @Description Verify OTP to activate user account and get access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Verify OTP Request"
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/user/verify-otp [post]
func (c *authController) VerifyOTP(ctx *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.authService.VerifyOTP(req)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Verification successful", dto.AuthResponse{Token: token})
}

// @Summary Login User
// @Description Login for users
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/user/login [post]
func (c *authController) LoginUser(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.authService.LoginUser(req)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Login successful", dto.AuthResponse{Token: token})
}

// @Summary Login Admin
// @Description Login for admins
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/admin/login [post]
func (c *authController) LoginAdmin(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.authService.LoginAdmin(req)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Login successful", dto.AuthResponse{Token: token})
}

// @Summary Login Worker
// @Description Login for workers
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/worker/login [post]
func (c *authController) LoginWorker(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.authService.LoginWorker(req)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Login successful", dto.AuthResponse{Token: token})
}

// @Summary Get Profile
// @Description Get current user's profile
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/me [get]
func (c *authController) GetProfile(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID := userIDVal.(uuid.UUID)
	profile, err := c.authService.GetProfile(userID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Profile retrieved", profile)
}

// @Summary Get All Users
// @Description Get all users (Admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserResponse
// @Failure 403 {object} map[string]string
// @Router /api/auth/admin/users [get]
func (c *authController) GetAllUsers(ctx *gin.Context) {
	users, err := c.authService.GetAllUsers()
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Users retrieved", users)
}

// @Summary Get All Workers
// @Description Get all workers (Admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserResponse
// @Failure 403 {object} map[string]string
// @Router /api/auth/admin/workers [get]
func (c *authController) GetAllWorkers(ctx *gin.Context) {
	workers, err := c.authService.GetAllWorkers()
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Workers retrieved", workers)
}

// @Summary Google Authentication
// @Description Authenticate with Google ID token and get JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.GoogleAuthRequest true "Google Auth Request"
// @Success 200 {object} dto.GoogleAuthResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/google [post]
func (c *authController) GoogleAuth(ctx *gin.Context) {
	var req dto.GoogleAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response, err := c.authService.GoogleAuth(req)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, "Authentication successful", response)
}
