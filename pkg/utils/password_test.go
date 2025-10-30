package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePasswordComplexity(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errCount int
	}{
		{
			name:     "Valid strong password",
			password: "StrongP@ssw0rd!",
			wantErr:  false,
			errCount: 0,
		},
		{
			name:     "Valid minimum requirements",
			password: "Secur3@!",
			wantErr:  false,
			errCount: 0,
		},
		{
			name:     "Too short password",
			password: "Test@1",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "No uppercase letter",
			password: "secur3@!",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "No lowercase letter",
			password: "SECUR3@!",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "No numbers",
			password: "SecurE@!",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "No special characters",
			password: "Secur3Pass",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "Contains forbidden word",
			password: "MyPassword@124",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "Multiple violations",
			password: "password",
			wantErr:  true,
			errCount: 4, // no uppercase, no numbers, no special, forbidden word
		},
		{
			name:     "Repeated characters",
			password: "Testaaa@123",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "Sequential numbers",
			password: "SecurE@012",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "Sequential letters",
			password: "Testabc@123",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "Too long password",
			password: "VeryLongPasswordThatExceedsTheMaximumLengthRequirementAndShouldFailValidationBecauseItIsTooLongForSecurityReasonsAndUsabilityP@ssw0rd!",
			wantErr:  true,
			errCount: 2, // too long, forbidden word
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordComplexity(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if complexityErr, ok := err.(*PasswordComplexityError); ok {
					assert.Len(t, complexityErr.Errors, tt.errCount, "Expected %d errors, got %d", tt.errCount, len(complexityErr.Errors))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePasswordComplexityCustomRequirements(t *testing.T) {
	customReq := PasswordComplexityRequirements{
		MinLength:        12,
		MaxLength:        50,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   false, // No special characters required
		ForbiddenWords:   []string{"test", "admin"},
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid with custom requirements",
			password: "MySecur3word",
			wantErr:  false,
		},
		{
			name:     "Too short for custom requirements",
			password: "Test123",
			wantErr:  true,
		},
		{
			name:     "Contains custom forbidden word",
			password: "AdminPassword123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordComplexity(tt.password, customReq)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAndHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password should be hashed",
			password: "StrongP@ssw0rd!",
			wantErr:  false,
		},
		{
			name:     "Invalid password should return validation error",
			password: "weak",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := ValidateAndHashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				assert.NotEqual(t, tt.password, hashedPassword)

				// Verify the hashed password can be checked
				assert.True(t, CheckPasswordHash(tt.password, hashedPassword))
			}
		})
	}
}

func TestDefaultPasswordRequirements(t *testing.T) {
	req := DefaultPasswordRequirements()

	assert.Equal(t, 8, req.MinLength)
	assert.Equal(t, 128, req.MaxLength)
	assert.True(t, req.RequireUppercase)
	assert.True(t, req.RequireLowercase)
	assert.True(t, req.RequireNumbers)
	assert.True(t, req.RequireSpecial)
	assert.Contains(t, req.ForbiddenWords, "password")
	assert.Contains(t, req.ForbiddenWords, "123456")
}

func TestPasswordComplexityError(t *testing.T) {
	err := &PasswordComplexityError{
		Message: "Test error",
		Errors:  []string{"Error 1", "Error 2"},
	}

	assert.Equal(t, "Test error", err.Error())
	assert.Len(t, err.Errors, 2)
}

// Benchmark tests for password validation performance
func BenchmarkValidatePasswordComplexity(b *testing.B) {
	password := "StrongP@ssw0rd123!"

	for i := 0; i < b.N; i++ {
		ValidatePasswordComplexity(password)
	}
}

func BenchmarkValidateAndHashPassword(b *testing.B) {
	password := "StrongP@ssw0rd123!"

	for i := 0; i < b.N; i++ {
		ValidateAndHashPassword(password)
	}
}
