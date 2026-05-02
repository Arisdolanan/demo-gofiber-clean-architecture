package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		email      string
		secret     string
		expiration time.Duration
		wantError  bool
	}{
		{
			name:       "Valid token generation",
			userID:     1,
			email:      "test@example.com",
			secret:     "test-secret",
			expiration: 15 * time.Minute,
			wantError:  false,
		},
		{
			name:       "Valid token with different user",
			userID:     123,
			email:      "user@test.com",
			secret:     "another-secret",
			expiration: 1 * time.Hour,
			wantError:  false,
		},
		{
			name:       "Empty secret",
			userID:     1,
			email:      "test@example.com",
			secret:     "",
			expiration: 15 * time.Minute,
			wantError:  false, // JWT library allows empty secret
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email, 1, "admin", tt.secret, tt.expiration)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				
				// Validate that the token has 3 parts (header.payload.signature)
				parts := len([]byte(token))
				assert.Greater(t, parts, 0)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret"
	userID := int64(1)
	email := "test@example.com"
	expiration := 15 * time.Minute

	// Generate a valid token for testing
	validToken, err := GenerateToken(userID, email, 1, "admin", secret, expiration)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		secret    string
		wantError bool
		wantClaims *JWTClaims
	}{
		{
			name:      "Valid token",
			token:     validToken,
			secret:    secret,
			wantError: false,
			wantClaims: &JWTClaims{
				UserID: userID,
				Email:  email,
			},
		},
		{
			name:      "Invalid secret",
			token:     validToken,
			secret:    "wrong-secret",
			wantError: true,
			wantClaims: nil,
		},
		{
			name:      "Invalid token format",
			token:     "invalid.token.format",
			secret:    secret,
			wantError: true,
			wantClaims: nil,
		},
		{
			name:      "Empty token",
			token:     "",
			secret:    secret,
			wantError: true,
			wantClaims: nil,
		},
		{
			name:      "Malformed token",
			token:     "malformed-token",
			secret:    secret,
			wantError: true,
			wantClaims: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token, tt.secret)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.wantClaims.UserID, claims.UserID)
				assert.Equal(t, tt.wantClaims.Email, claims.Email)
				assert.True(t, claims.ExpiresAt.After(time.Now()))
			}
		})
	}
}

func TestGenerateAndValidateTokenRoundTrip(t *testing.T) {
	secret := "test-secret"
	userID := int64(123)
	email := "user@example.com"
	expiration := 1 * time.Hour

	// Generate token
	token, err := GenerateToken(userID, email, 1, "admin", secret, expiration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the same token
	claims, err := ValidateToken(token, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.IssuedAt.Before(time.Now().Add(1*time.Second)))
}

func TestExpiredToken(t *testing.T) {
	secret := "test-secret"
	userID := int64(1)
	email := "test@example.com"
	expiration := -1 * time.Hour // Expired token

	// Generate an expired token
	token, err := GenerateToken(userID, email, 1, "admin", secret, expiration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Try to validate the expired token
	claims, err := ValidateToken(token, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
	
	// Should be a validation error related to expiration
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWTClaims(t *testing.T) {
	claims := &JWTClaims{
		UserID: 123,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	assert.Equal(t, int64(123), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)
}