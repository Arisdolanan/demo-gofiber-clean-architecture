package config

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupLimiter() limiter.Config {
	return limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}
}
