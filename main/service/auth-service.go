package service

import (
	"errors"
	"fmt"

	"github.com/inlovewithgo/transit-backend/main/models"
	repo "github.com/inlovewithgo/transit-backend/main/repo/interface"
	"github.com/inlovewithgo/transit-backend/main/utils"
	"github.com/inlovewithgo/transit-backend/pkg/logger"
)

type AuthService struct {
	userRepo    repo.UserRepository
	mailService *MailService
}

func NewAuthService(userRepo repo.UserRepository, mailService *MailService) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		mailService: mailService,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	exists, err := s.userRepo.UserExists(req.Email)
	if err != nil {
		logger.Log.Error("Error checking user existence: %v", err)
		return nil, fmt.Errorf("failed to check user existence")
	}

	if exists {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error("Error hashing password: %v", err)
		return nil, fmt.Errorf("failed to hash password")
	}

	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		logger.Log.Error("Error creating user: %v", err)
		return nil, fmt.Errorf("failed to create user")
	}

	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		logger.Log.Error("Error generating access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token")
	}

	go func() {
		err := s.mailService.SendWelcomeEmail(user.Email, user.FirstName, user.LastName)
		if err != nil {
			logger.Log.Error("Failed to send welcome email: %v", err)
		}
	}()

	user.Password = ""

	return &models.AuthResponse{
		User:        *user,
		AccessToken: accessToken,
		Message:     "User registered successfully",
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		logger.Log.Error("Error getting user by email: %v", err)
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		logger.Log.Error("Error generating access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token")
	}

	go func() {
		err := s.mailService.SendLoginNotification(user.Email, user.FirstName, user.LastName)
		if err != nil {
			logger.Log.Error("Failed to send login notification email: %v", err)
		}
	}()

	user.Password = ""

	return &models.AuthResponse{
		User:        *user,
		AccessToken: accessToken,
		Message:     "Login successful",
	}, nil
}

func (s *AuthService) GetUserProfile(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		logger.Log.Error("Error getting user profile: %v", err)
		return nil, errors.New("user not found")
	}

	user.Password = ""
	return user, nil
}
