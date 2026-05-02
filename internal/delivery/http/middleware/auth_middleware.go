package middleware

import (
	"log"
	"strings"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/redis"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// JWTProtected middleware for JWT authentication with token blacklisting
func JWTProtected(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Authorization header is required",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Invalid authorization header format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Token is required",
			})
		}

		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Invalid or expired token",
			})
		}

		// Set user information and token in context for blacklist checking
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("school_id", claims.SchoolID)
		c.Locals("user_type", claims.UserType)
		c.Locals("token", tokenString)

		return c.Next()
	}
}

// JWTProtectedWithBlacklist middleware for JWT authentication with Redis blacklist checking
func JWTProtectedWithBlacklist(jwtSecret string, authRedis redis.AuthRedisRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Authorization header is required",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Invalid authorization header format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Token is required",
			})
		}

		// Check if token is blacklisted (gracefully handle Redis unavailability)
		isBlacklisted := false
		isBlacklisted, err := authRedis.IsTokenBlacklisted(tokenString)
		if err != nil {
			// Log the error but continue - Redis failure shouldn't block all requests
			log.Printf("Warning: Redis unavailable, skipping blacklist check: %v", err)
		}

		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Token has been revoked",
			})
		}

		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("school_id", claims.SchoolID)
		c.Locals("user_type", claims.UserType)
		c.Locals("token", tokenString)

		return c.Next()
	}
}
