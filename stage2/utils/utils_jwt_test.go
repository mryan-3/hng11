package utils

import (
	"os"
	"testing"
	"time"

	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestSignJwtToken(t *testing.T) {
	// Set up the environment variable for JWT_SECRET
	os.Setenv("JWT_SECRET", "your-secret-key")

	// Test signing a JWT token
	userID := "testUserID"
	tokenString, err := SignJwtToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify its content
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	assert.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["user_id"])
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	assert.WithinDuration(t, time.Now().Add(time.Hour*24*7), expirationTime, time.Minute)
}



func TestVerifyJwtToken(t *testing.T) {
	// Set up the environment variable for JWT_SECRET
	os.Setenv("JWT_SECRET", "your-secret-key")

	// Create a token for testing
	userID := "testUserID"
	tokenString, err := SignJwtToken(userID)
	assert.NoError(t, err)

	// Verify the token
	claims, isValid, err := VerifyJwtToken(tokenString)
	assert.NoError(t, err)
	assert.True(t, isValid)
	assert.Equal(t, userID, claims["user_id"])
}

func TestVerifyInvalidJwtToken(t *testing.T) {
	// Set up the environment variable for JWT_SECRET
	os.Setenv("JWT_SECRET", "your-secret-key")

	// Use an invalid token
	invalidTokenString := "invalid.token.string"

	// Verify the token
	claims, isValid, err := VerifyJwtToken(invalidTokenString)
	assert.Error(t, err)
	assert.False(t, isValid)
	assert.Nil(t, claims)
}
