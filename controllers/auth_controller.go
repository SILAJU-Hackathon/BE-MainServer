package controllers

import (
	"net/http"

	"dinacom-11.0-backend/models/dto"
	"dinacom-11.0-backend/services"
	"dinacom-11.0-backend/utils"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	RegisterUser(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	LoginAdmin(ctx *gin.Context)
	LoginWorker(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService: authService}
}

// RegisterUser godoc
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

// VerifyOTP godoc
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

// LoginUser godoc
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

// LoginAdmin godoc
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

// LoginWorker godoc
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
