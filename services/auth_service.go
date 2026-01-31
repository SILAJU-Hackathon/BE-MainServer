package services

import (
	"errors"

	"dinacom-11.0-backend/models/dto"
	entity "dinacom-11.0-backend/models/entity"
	"dinacom-11.0-backend/repositories"
	"dinacom-11.0-backend/utils"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(req dto.RegisterRequest) error
	LoginUser(req dto.LoginRequest) (string, error)
	LoginAdmin(req dto.LoginRequest) (string, error)
	LoginWorker(req dto.LoginRequest) (string, error)
	GoogleAuth(req dto.GoogleAuthRequest) (*dto.GoogleAuthResponse, error)
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	GetAllUsers() ([]dto.UserResponse, error)
	GetAllWorkers() ([]dto.UserResponse, error)
	AssignWorkerRole(userID uuid.UUID) error
	CreateWorker(req dto.CreateWorkerRequest) (*dto.WorkerResponse, error)
	GetWorkerByID(workerID uuid.UUID) (*dto.WorkerResponse, error)
	UpdateWorker(workerID uuid.UUID, req dto.UpdateWorkerRequest) (*dto.WorkerResponse, error)
	DeleteWorker(workerID uuid.UUID) error
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) RegisterUser(req dto.RegisterRequest) error {
	existingUser, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return err
	}

	// If user already exists, reject registration
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Create new user (verified by default, no OTP needed)
	user := &entity.User{
		Username: req.Username,
		Fullname: req.FullName,
		Email:    req.Email,
		Role:     "user",
		Password: hashedPassword,
		Verified: true,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

// VerifyOTP is deprecated - users are now verified on registration
func (s *authService) VerifyOTP(req dto.VerifyOTPRequest) (string, error) {
	return "", errors.New("OTP verification is no longer required")
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

func (s *authService) CreateWorker(req dto.CreateWorkerRequest) (*dto.WorkerResponse, error) {
	existingUser, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	worker := &entity.User{
		Username: req.Username,
		Fullname: req.Fullname,
		Email:    req.Email,
		Role:     entity.ROLE_WORKER,
		Password: hashedPassword,
		Verified: true,
	}

	if err := s.userRepo.CreateUser(worker); err != nil {
		return nil, err
	}

	return &dto.WorkerResponse{
		ID:        worker.ID,
		Username:  worker.Username,
		Fullname:  worker.Fullname,
		Email:     worker.Email,
		Verified:  worker.Verified,
		CreatedAt: worker.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *authService) GetWorkerByID(workerID uuid.UUID) (*dto.WorkerResponse, error) {
	worker, err := s.userRepo.FindUserByID(workerID)
	if err != nil {
		return nil, err
	}
	if worker == nil {
		return nil, errors.New("worker not found")
	}
	if worker.Role != entity.ROLE_WORKER {
		return nil, errors.New("user is not a worker")
	}

	return &dto.WorkerResponse{
		ID:        worker.ID,
		Username:  worker.Username,
		Fullname:  worker.Fullname,
		Email:     worker.Email,
		Verified:  worker.Verified,
		CreatedAt: worker.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *authService) UpdateWorker(workerID uuid.UUID, req dto.UpdateWorkerRequest) (*dto.WorkerResponse, error) {
	worker, err := s.userRepo.FindUserByID(workerID)
	if err != nil {
		return nil, err
	}
	if worker == nil {
		return nil, errors.New("worker not found")
	}
	if worker.Role != entity.ROLE_WORKER {
		return nil, errors.New("user is not a worker")
	}

	if req.Fullname != "" {
		worker.Fullname = req.Fullname
	}
	if req.Username != "" {
		worker.Username = req.Username
	}
	if req.Email != "" {
		worker.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		worker.Password = hashedPassword
	}

	if err := s.userRepo.UpdateUser(worker); err != nil {
		return nil, err
	}

	return &dto.WorkerResponse{
		ID:        worker.ID,
		Username:  worker.Username,
		Fullname:  worker.Fullname,
		Email:     worker.Email,
		Verified:  worker.Verified,
		CreatedAt: worker.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *authService) DeleteWorker(workerID uuid.UUID) error {
	worker, err := s.userRepo.FindUserByID(workerID)
	if err != nil {
		return err
	}
	if worker == nil {
		return errors.New("worker not found")
	}
	if worker.Role != entity.ROLE_WORKER {
		return errors.New("user is not a worker")
	}

	return s.userRepo.DeleteUser(workerID)
}
