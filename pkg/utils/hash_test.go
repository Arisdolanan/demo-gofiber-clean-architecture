package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		{
			name:      "Valid password",
			password:  "password123",
			wantError: false,
		},
		{
			name:      "Empty password",
			password:  "",
			wantError: false, // bcrypt allows empty passwords
		},
		{
			name:      "Long password (over 72 bytes)",
			password:  "this-is-a-very-long-password-that-should-still-work-fine-with-bcrypt-hashing-algorithm",
			wantError: true,
		},
		{
			name:      "Long password (under 72 bytes)",
			password:  "this-is-a-long-password-but-under-72-bytes",
			wantError: false,
		},
		{
			name:      "Password with special characters",
			password:  "p@ssw0rd!@#$%^&*()",
			wantError: false,
		},
		{
			name:      "Unicode password",
			password:  "пароль123",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash) // Hash should be different from password
				assert.True(t, len(hash) > 0)
				
				// Verify that the hash starts with bcrypt identifier
				if len(hash) > 4 {
					assert.Equal(t, "$2a$", hash[:4]) // bcrypt format identifier
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			expected: true,
		},
		{
			name:     "Wrong password",
			password: "wrongpassword",
			hash:     hash,
			expected: false,
		},
		{
			name:     "Empty password with valid hash",
			password: "",
			hash:     hash,
			expected: false,
		},
		{
			name:     "Valid password with empty hash",
			password: password,
			hash:     "",
			expected: false,
		},
		{
			name:     "Empty password and empty hash",
			password: "",
			hash:     "",
			expected: false,
		},
		{
			name:     "Invalid hash format",
			password: password,
			hash:     "invalid-hash",
			expected: false,
		},
		{
			name:     "Case sensitive password",
			password: "TESTPASSWORD123",
			hash:     hash,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPasswordHash(tt.password, tt.hash)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "testpassword"
	
	// Hash the same password multiple times
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)
	hash3, err3 := HashPassword(password)
	
	// All should succeed
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	
	// All hashes should be different (bcrypt includes random salt)
	assert.NotEqual(t, hash1, hash2)
	assert.NotEqual(t, hash2, hash3)
	assert.NotEqual(t, hash1, hash3)
	
	// But all should validate correctly
	assert.True(t, CheckPasswordHash(password, hash1))
	assert.True(t, CheckPasswordHash(password, hash2))
	assert.True(t, CheckPasswordHash(password, hash3))
}

func TestHashPasswordAndCheckRoundTrip(t *testing.T) {
	testCases := []string{
		"simplepassword",
		"complex!P@ssw0rd#123",
		"",
		"a",
		"medium-length-password-that-works-with-bcrypt",
		"🔐🔑password🔒",
	}
	
	for _, password := range testCases {
		t.Run("password: "+password, func(t *testing.T) {
			// Skip if password is too long for bcrypt
			if len(password) > 72 {
				return
			}
					
			// Hash the password
			hash, err := HashPassword(password)
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
			
			// Check that the password validates against its hash
			assert.True(t, CheckPasswordHash(password, hash))
			
			// Check that a wrong password doesn't validate
			assert.False(t, CheckPasswordHash(password+"wrong", hash))
		})
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	password := "benchmarkpassword123"
	hash, _ := HashPassword(password)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CheckPasswordHash(password, hash)
	}
}