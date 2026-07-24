package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret-key"

func TestGenerateAndValidateToken_ValidRoundTrip(t *testing.T) {
	token, err := GenerateToken("user123", "john@example.com", "user", testSecret, 15*time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token, testSecret)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "john@example.com", claims.Email)
	assert.Equal(t, "user", claims.Role)
}

func TestValidateToken_ExpiredToken_ReturnsExpiredError(t *testing.T) {
	// Generate a token that expired 1 minute ago
	token, err := GenerateToken("user123", "john@example.com", "user", testSecret, -1*time.Minute)
	assert.NoError(t, err)

	_, err = ValidateToken(token, testSecret)
	assert.ErrorIs(t, err, ErrExpiredToken)
}

func TestValidateToken_MalformedToken_ReturnsInvalidError(t *testing.T) {
	_, err := ValidateToken("this.is.not-a-valid-jwt", testSecret)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestValidateToken_WrongSecret_ReturnsInvalidError(t *testing.T) {
	token, err := GenerateToken("user123", "john@example.com", "user", testSecret, 15*time.Minute)
	assert.NoError(t, err)

	_, err = ValidateToken(token, "a-completely-different-secret")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestValidateToken_TamperedSignature_ReturnsInvalidError(t *testing.T) {
	token, err := GenerateToken("user123", "john@example.com", "user", testSecret, 15*time.Minute)
	assert.NoError(t, err)

	// Flip the last character of the token — corrupts the signature
	tampered := token[:len(token)-1] + "x"

	_, err = ValidateToken(tampered, testSecret)
	assert.ErrorIs(t, err, ErrInvalidToken)
}