package service

import (
	"context"
	"testing"
	"time"

	"github.com/ararext/Go-JWT-Authentication-API/internal/dto"
	"github.com/ararext/Go-JWT-Authentication-API/internal/repository/mock"
	"github.com/stretchr/testify/assert"
)

func newTestAuthService() AuthService {
	repo := mock.NewUserRepository()
	return NewAuthService(repo, "test-secret", 15*time.Minute, 7*24*time.Hour)
}

func TestSignup_Success(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	resp, err := svc.Signup(ctx, dto.SignupRequest{
		Name:     "Ararext",
		Email:    "ararext@example.com",
		Password: "securepass123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, "ararext@example.com", resp.User.Email)
	assert.Equal(t, "user", resp.User.Role)
}

func TestSignup_DuplicateEmail_ReturnsError(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	req := dto.SignupRequest{
		Name:     "Ararext",
		Email:    "ararext@example.com",
		Password: "securepass123",
	}

	_, err := svc.Signup(ctx, req)
	assert.NoError(t, err)

	// Second signup with same email should fail
	_, err = svc.Signup(ctx, req)
	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestLogin_Success(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	_, err := svc.Signup(ctx, dto.SignupRequest{
		Name:     "Ararext",
		Email:    "ararext@example.com",
		Password: "securepass123",
	})
	assert.NoError(t, err)

	resp, err := svc.Login(ctx, dto.LoginRequest{
		Email:    "ararext@example.com",
		Password: "securepass123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestLogin_WrongPassword_ReturnsInvalidCredentials(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	_, err := svc.Signup(ctx, dto.SignupRequest{
		Name:     "Ararext",
		Email:    "ararext@example.com",
		Password: "securepass123",
	})
	assert.NoError(t, err)

	_, err = svc.Login(ctx, dto.LoginRequest{
		Email:    "ararext@example.com",
		Password: "wrongpassword",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_UserNotFound_ReturnsInvalidCredentials(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	_, err := svc.Login(ctx, dto.LoginRequest{
		Email:    "doesnotexist@example.com",
		Password: "whatever123",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}