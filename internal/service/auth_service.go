package service

import (
	"context"
	"errors"
	"time"

	"github.com/ararext/Go-JWT-Authentication-API/internal/dto"
	"github.com/ararext/Go-JWT-Authentication-API/internal/models"
	"github.com/ararext/Go-JWT-Authentication-API/internal/repository"
	"github.com/ararext/Go-JWT-Authentication-API/internal/utils"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	Signup(ctx context.Context, req dto.SignupRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo             repository.UserRepository
	jwtSecret            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtSecret string,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) AuthService {
	return &authService{
		userRepo:             userRepo,
		jwtSecret:            jwtSecret,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (s *authService) Signup(ctx context.Context, req dto.SignupRequest) (*dto.AuthResponse, error) {
	// Friendly pre-check — the DB-level unique index from Day 2 is the real guarantee,
	// this just avoids a generic 500 error on the common case.
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err // some real DB error, not "not found"
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(user)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.buildAuthResponse(user)
}

func (s *authService) buildAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	userID := user.ID.Hex()

	accessToken, err := utils.GenerateToken(userID, user.Email, user.Role, s.jwtSecret, s.accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(userID, user.Email, user.Role, s.jwtSecret, s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    userID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}