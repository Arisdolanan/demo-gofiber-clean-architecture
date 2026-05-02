package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	SchoolID int64  `json:"school_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token
func GenerateToken(userID int64, email string, schoolID int64, userType string, secret string, expiration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		SchoolID: schoolID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GetUserIDFromToken extracts user ID from JWT token in Fiber context
func GetUserIDFromToken(ctx *fiber.Ctx) (int64, error) {
	authorization := ctx.Get("Authorization")
	if authorization == "" {
		return 0, errors.New("authorization header missing")
	}

	parts := strings.Split(authorization, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, errors.New("invalid authorization header format")
	}

	tokenString := parts[1]
	if tokenString == "" {
		return 0, errors.New("token missing")
	}

	if userID := ctx.Locals("user_id"); userID != nil {
		if id, ok := userID.(int64); ok {
			return id, nil
		}
		if id, ok := userID.(float64); ok {
			return int64(id), nil
		}
	}

	return 0, errors.New("user ID not found in context")
}

// GetSchoolIDFromToken extracts school ID from JWT token in Fiber context
func GetSchoolIDFromToken(ctx *fiber.Ctx) (int64, error) {
	if schoolID := ctx.Locals("school_id"); schoolID != nil {
		if id, ok := schoolID.(int64); ok {
			return id, nil
		}
		if id, ok := schoolID.(float64); ok {
			return int64(id), nil
		}
	}

	return 0, errors.New("school ID not found in context")
}
