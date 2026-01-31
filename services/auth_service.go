package services

import (
	"errors"
	"sync"

	"dinacom-11.0-backend/models/dto"
	entity "dinacom-11.0-backend/models/entity"
	"dinacom-11.0-backend/repositories"
	"dinacom-11.0-backend/utils"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(req dto.RegisterRequest) error
	VerifyOTP(req dto.VerifyOTPRequest) (string, error)
	LoginUser(req dto.LoginRequest) (string, error)
	LoginAdmin(req dto.LoginRequest) (string, error)
	LoginWorker(req dto.LoginRequest) (string, error)
	GoogleAuth(req dto.GoogleAuthRequest) (*dto.GoogleAuthResponse, error)
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	GetAllUsers() ([]dto.UserResponse, error)
	GetAllWorkers() ([]dto.UserResponse, error)
	AssignWorkerRole(userID uuid.UUID) error
}

type authService struct {
	userRepo repositories.UserRepository
	otpStore map[string]string
	mutex    sync.RWMutex
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		otpStore: make(map[string]string),
	}
}

func (s *authService) RegisterUser(req dto.RegisterRequest) error {
	existingUser, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

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

	if err := s.userRepo.UpdateUserVerified(req.Email, true); err != nil {
		return "", err
	}

	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", err
	}

	s.mutex.Lock()
	delete(s.otpStore, req.Email)
	s.mutex.Unlock()

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
		return "", errors.New("account not verified. please verify OTP")
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateAccessToken(user.ID, user.Role, user.Email)
}

func (s *authService) LoginAdmin(req dto.LoginRequest) (string, error) {
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

func (s *authService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Fullname: user.Fullname,
		Email:    user.Email,
		Role:     user.Role,
		Verified: user.Verified,
	}, nil
}

func (s *authService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetUsersByRole("user")
	if err != nil {
		return nil, err
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Fullname: user.Fullname,
			Email:    user.Email,
			Role:     user.Role,
			Verified: user.Verified,
		})
	}
	return response, nil
}

func (s *authService) GetAllWorkers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetUsersByRole("worker")
	if err != nil {
		return nil, err
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Fullname: user.Fullname,
			Email:    user.Email,
			Role:     user.Role,
			Verified: user.Verified,
		})
	}
	return response, nil
}

func (s *authService) GoogleAuth(req dto.GoogleAuthRequest) (*dto.GoogleAuthResponse, error) {
	tokenInfo, err := utils.VerifyGoogleToken(req.IDToken)
	if err != nil {
		return nil, errors.New("invalid Google token")
	}

	existingUser, _ := s.userRepo.FindUserByEmail(tokenInfo.Email)
	isNewUser := false

	if existingUser == nil {
		isNewUser = true
		newUser := &entity.User{
			Username: tokenInfo.Email,
			Fullname: tokenInfo.Name,
			Email:    tokenInfo.Email,
			Role:     "user",
			Password: "",
			Verified: true,
		}
		if err := s.userRepo.CreateUser(newUser); err != nil {
			return nil, err
		}
		existingUser, _ = s.userRepo.FindUserByEmail(tokenInfo.Email)
	}

	token, err := utils.GenerateAccessToken(existingUser.ID, existingUser.Role, existingUser.Email)
	if err != nil {
		return nil, err
	}

	return &dto.GoogleAuthResponse{
		Token: token,
		User: dto.GoogleUserDetails{
			Email:     existingUser.Email,
			Fullname:  existingUser.Fullname,
			IsNewUser: isNewUser,
		},
	}, nil
}

func (s *authService) AssignWorkerRole(userID uuid.UUID) error {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Role = entity.ROLE_WORKER
	return s.userRepo.UpdateUser(user)
}
