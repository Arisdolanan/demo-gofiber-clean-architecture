package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRedisRepository interface {
	StoreRefreshToken(userID int64, refreshToken string, expiration time.Duration) error
	GetRefreshToken(userID int64) (string, error)
	DeleteRefreshToken(userID int64) error
	// Token blacklisting methods
	BlacklistToken(tokenString string, expiration time.Duration) error
	IsTokenBlacklisted(tokenString string) (bool, error)
	BlacklistAllUserTokens(userID int64, expiration time.Duration) error
}

type authRedisRepository struct {
	client *redis.Client
}

func NewAuthRedisRepository(client *redis.Client) AuthRedisRepository {
	return &authRedisRepository{client: client}
}

// StoreRefreshToken stores a refresh token in Redis
func (r *authRedisRepository) StoreRefreshToken(userID int64, refreshToken string, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", userID)
	return r.client.Set(ctx, key, refreshToken, expiration).Err()
}

// GetRefreshToken retrieves a refresh token from Redis
func (r *authRedisRepository) GetRefreshToken(userID int64) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", userID)
	return r.client.Get(ctx, key).Result()
}

// DeleteRefreshToken removes a refresh token from Redis
func (r *authRedisRepository) DeleteRefreshToken(userID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", userID)
	return r.client.Del(ctx, key).Err()
}

// BlacklistToken adds a token to the blacklist
func (r *authRedisRepository) BlacklistToken(tokenString string, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist_token:%s", tokenString)
	return r.client.Set(ctx, key, "blacklisted", expiration).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted
func (r *authRedisRepository) IsTokenBlacklisted(tokenString string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist_token:%s", tokenString)
	_, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// BlacklistAllUserTokens blacklists all tokens for a specific user (used during logout)
func (r *authRedisRepository) BlacklistAllUserTokens(userID int64, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist_user:%d", userID)
	return r.client.Set(ctx, key, "blacklisted", expiration).Err()
}
