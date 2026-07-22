package utils

import (
	"testing"

	"github.com/ararext/Go-JWT-Authentication-API/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestValidateStruct_ValidSignupRequest_ReturnsNoError(t *testing.T) {
	req := dto.SignupRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	err := ValidateStruct(req)
	assert.NoError(t, err)
}

func TestValidateStruct_ShortPassword_ReturnsError(t *testing.T) {
	req := dto.SignupRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "short",
	}

	err := ValidateStruct(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Password")
}

func TestValidateStruct_InvalidEmail_ReturnsError(t *testing.T) {
	req := dto.SignupRequest{
		Name:     "John Doe",
		Email:    "not-an-email",
		Password: "password123",
	}

	err := ValidateStruct(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Email")
}