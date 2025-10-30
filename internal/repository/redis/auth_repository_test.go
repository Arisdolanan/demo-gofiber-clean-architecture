package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthRedisRepositoryTestSuite struct {
	suite.Suite
	repo      AuthRedisRepository
	client    *redis.Client
	miniRedis *miniredis.Miniredis
}

func (suite *AuthRedisRepositoryTestSuite) SetupTest() {
	// Start miniredis server
	mr, err := miniredis.Run()
	suite.Require().NoError(err)
	suite.miniRedis = mr

	// Create Redis client
	suite.client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Create repository
	suite.repo = NewAuthRedisRepository(suite.client)
}

func (suite *AuthRedisRepositoryTestSuite) TearDownTest() {
	suite.client.Close()
	suite.miniRedis.Close()
}

// Test StoreRefreshToken
func (suite *AuthRedisRepositoryTestSuite) TestStoreRefreshToken_Success() {
	userID := int64(1)
	refreshToken := "test-refresh-token"
	expiration := 24 * time.Hour

	err := suite.repo.StoreRefreshToken(userID, refreshToken, expiration)
	assert.NoError(suite.T(), err)

	// Verify token is stored
	storedToken, err := suite.repo.GetRefreshToken(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), refreshToken, storedToken)
}

// Test GetRefreshToken
func (suite *AuthRedisRepositoryTestSuite) TestGetRefreshToken_Success() {
	userID := int64(1)
	refreshToken := "test-refresh-token"
	expiration := 24 * time.Hour

	// Store token first
	err := suite.repo.StoreRefreshToken(userID, refreshToken, expiration)
	assert.NoError(suite.T(), err)

	// Get token
	storedToken, err := suite.repo.GetRefreshToken(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), refreshToken, storedToken)
}

func (suite *AuthRedisRepositoryTestSuite) TestGetRefreshToken_NotFound() {
	userID := int64(999)

	// Try to get non-existent token
	_, err := suite.repo.GetRefreshToken(userID)
	assert.Error(suite.T(), err)
}

// Test DeleteRefreshToken
func (suite *AuthRedisRepositoryTestSuite) TestDeleteRefreshToken_Success() {
	userID := int64(1)
	refreshToken := "test-refresh-token"
	expiration := 24 * time.Hour

	// Store token first
	err := suite.repo.StoreRefreshToken(userID, refreshToken, expiration)
	assert.NoError(suite.T(), err)

	// Delete token
	err = suite.repo.DeleteRefreshToken(userID)
	assert.NoError(suite.T(), err)

	// Verify token is deleted
	_, err = suite.repo.GetRefreshToken(userID)
	assert.Error(suite.T(), err)
}

// Test BlacklistToken
func (suite *AuthRedisRepositoryTestSuite) TestBlacklistToken_Success() {
	tokenString := "test-access-token"
	expiration := 15 * time.Minute

	err := suite.repo.BlacklistToken(tokenString, expiration)
	assert.NoError(suite.T(), err)

	// Verify token is blacklisted
	isBlacklisted, err := suite.repo.IsTokenBlacklisted(tokenString)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isBlacklisted)
}

// Test IsTokenBlacklisted
func (suite *AuthRedisRepositoryTestSuite) TestIsTokenBlacklisted_Success() {
	tokenString := "test-access-token"
	expiration := 15 * time.Minute

	// Token should not be blacklisted initially
	isBlacklisted, err := suite.repo.IsTokenBlacklisted(tokenString)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isBlacklisted)

	// Blacklist token
	err = suite.repo.BlacklistToken(tokenString, expiration)
	assert.NoError(suite.T(), err)

	// Now token should be blacklisted
	isBlacklisted, err = suite.repo.IsTokenBlacklisted(tokenString)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isBlacklisted)
}

func (suite *AuthRedisRepositoryTestSuite) TestIsTokenBlacklisted_NotBlacklisted() {
	tokenString := "non-blacklisted-token"

	isBlacklisted, err := suite.repo.IsTokenBlacklisted(tokenString)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isBlacklisted)
}

// Test BlacklistAllUserTokens
func (suite *AuthRedisRepositoryTestSuite) TestBlacklistAllUserTokens_Success() {
	userID := int64(1)
	expiration := 24 * time.Hour

	err := suite.repo.BlacklistAllUserTokens(userID, expiration)
	assert.NoError(suite.T(), err)

	// Verify the blacklist entry exists
	// Note: This is testing the storage, in practice you'd check against this in middleware
	key := "blacklist_user:1"
	exists := suite.miniRedis.Exists(key)
	assert.True(suite.T(), exists)
}

// Test token expiration
func (suite *AuthRedisRepositoryTestSuite) TestTokenExpiration() {
	tokenString := "expiring-token"
	shortExpiration := 1 * time.Millisecond

	// Blacklist token with very short expiration
	err := suite.repo.BlacklistToken(tokenString, shortExpiration)
	assert.NoError(suite.T(), err)

	// Wait for expiration
	time.Sleep(2 * time.Millisecond)

	// Fast forward time in miniredis
	suite.miniRedis.FastForward(shortExpiration)

	// Token should no longer be blacklisted
	isBlacklisted, err := suite.repo.IsTokenBlacklisted(tokenString)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isBlacklisted)
}

// Test concurrent operations
func (suite *AuthRedisRepositoryTestSuite) TestConcurrentOperations() {
	userID := int64(1)
	refreshToken := "concurrent-test-token"
	expiration := 24 * time.Hour

	// Run concurrent operations
	done := make(chan bool, 2)

	go func() {
		err := suite.repo.StoreRefreshToken(userID, refreshToken, expiration)
		assert.NoError(suite.T(), err)
		done <- true
	}()

	go func() {
		time.Sleep(10 * time.Millisecond) // Small delay
		_, _ = suite.repo.GetRefreshToken(userID)
		// Error is acceptable here due to timing
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done
}

// Integration test for full token lifecycle
func (suite *AuthRedisRepositoryTestSuite) TestTokenLifecycle() {
	userID := int64(1)
	refreshToken := "lifecycle-test-token"
	accessToken := "lifecycle-access-token"
	expiration := 24 * time.Hour

	// 1. Store refresh token
	err := suite.repo.StoreRefreshToken(userID, refreshToken, expiration)
	assert.NoError(suite.T(), err)

	// 2. Verify refresh token exists
	storedToken, err := suite.repo.GetRefreshToken(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), refreshToken, storedToken)

	// 3. Blacklist access token (simulate logout)
	err = suite.repo.BlacklistToken(accessToken, expiration)
	assert.NoError(suite.T(), err)

	// 4. Verify access token is blacklisted
	isBlacklisted, err := suite.repo.IsTokenBlacklisted(accessToken)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isBlacklisted)

	// 5. Delete refresh token (complete logout)
	err = suite.repo.DeleteRefreshToken(userID)
	assert.NoError(suite.T(), err)

	// 6. Verify refresh token is deleted
	_, err = suite.repo.GetRefreshToken(userID)
	assert.Error(suite.T(), err)
}

func TestAuthRedisRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRedisRepositoryTestSuite))
}
