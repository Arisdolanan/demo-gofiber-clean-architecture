package config

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupLimiter() limiter.Config {
	return limiter.Config{
		Max:               100,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	}
}
