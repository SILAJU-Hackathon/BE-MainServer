package services

import (
	"errors"
	"sync"

	"dinacom-11.0-backend/models/dto"
	entity "dinacom-11.0-backend/models/entity"
	"dinacom-11.0-backend/repositories"
	"dinacom-11.0-backend/utils"
)

type AuthService interface {
	RegisterUser(req dto.RegisterRequest) error
	VerifyOTP(req dto.VerifyOTPRequest) (string, error)
	LoginUser(req dto.LoginRequest) (string, error)
	LoginAdmin(req dto.LoginRequest) (string, error)
	LoginWorker(req dto.LoginRequest) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
	otpStore map[string]string // Simple in-memory store for OTPs: email -> otp
	mutex    sync.RWMutex
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		otpStore: make(map[string]string),
	}
}

func (s *authService) RegisterUser(req dto.RegisterRequest) error {
	// Check if user exists
	existingUser, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Create user (verified = false)
	user := &entity.User{
		Username: req.Username,
		Fullname: req.FullName,
		Email:    req.Email,
		Role:     "user",
		Password: hashedPassword,
		Verified: false,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return err
	}

	// Generate and Send OTP
	otp := utils.GenerateOTP()

	s.mutex.Lock()
	s.otpStore[req.Email] = otp
	s.mutex.Unlock()

	utils.SendOTP(req.Email, otp) 

	return nil
}

func (s *authService) VerifyOTP(req dto.VerifyOTPRequest) (string, error) {
	s.mutex.RLock()
	storedOTP, exists := s.otpStore[req.Email]
	s.mutex.RUnlock()

	if !exists || storedOTP != req.OTP {
		return "", errors.New("invalid or expired OTP")
	}

	// Verify user in DB
	if err := s.userRepo.UpdateUserVerified(req.Email, true); err != nil {
		return "", err
	}

	// Get User to generate token
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", err
	}

	// Clear OTP
	s.mutex.Lock()
	delete(s.otpStore, req.Email)
	s.mutex.Unlock()

	// Generate Token
	return utils.GenerateAccessToken(user.ID, user.Role, user.Email)
}

func (s *authService) LoginUser(req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	if user.Role != "user" {
		return "", errors.New("unauthorized: user role required")
	}

	if !user.Verified {
		// Prepare new OTP if not verified (optional, but good UX to allow re-verify)
		// For now, just block
		return "", errors.New("account not verified. please verify OTP")
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateAccessToken(user.ID, user.Role, user.Email)
}

func (s *authService) LoginAdmin(req dto.LoginRequest) (string, error) {
	// For admin, we assume they might be pre-seeded or created via different flow.
	// But validation logic is similar, just role check.
	// NOTE: If Admin doesn't exist, we can't login.
	// For testing, user might need to seed an admin.

	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	if user.Role != "admin" {
		return "", errors.New("unauthorized: admin role required")
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateAccessToken(user.ID, user.Role, user.Email)
}

func (s *authService) LoginWorker(req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	if user.Role != "worker" {
		return "", errors.New("unauthorized: worker role required")
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateAccessToken(user.ID, user.Role, user.Email)
}
