package passwordutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	rawPassword := "12345678"
	hashed, err := HashPassword(rawPassword)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, rawPassword, hashed)
}

func TestPasswordMatches(t *testing.T) {
	rawPassword := "12345678"
	hashed, _ := HashPassword(rawPassword)

	equals := PasswordMatches(hashed, rawPassword)

	assert.True(t, equals)
}
