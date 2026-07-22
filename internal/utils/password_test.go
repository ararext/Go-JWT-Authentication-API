package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_ProducesDifferentHashForSamePassword(t *testing.T) {
	password := "mySecret123"

	hash1, err1 := HashPassword(password)
	assert.NoError(t, err1)
	assert.NotEmpty(t, hash1)

	hash2, err2 := HashPassword(password)
	assert.NoError(t, err2)

	// Same input, different hash — proves bcrypt is salting correctly
	assert.NotEqual(t, hash1, hash2)
}

func TestCheckPassword_CorrectPassword_ReturnsTrue(t *testing.T) {
	password := "mySecret123"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	result := CheckPassword(password, hash)
	assert.True(t, result)
}

func TestCheckPassword_WrongPassword_ReturnsFalse(t *testing.T) {
	hash, err := HashPassword("correctPassword")
	assert.NoError(t, err)

	result := CheckPassword("wrongPassword", hash)
	assert.False(t, result)
}