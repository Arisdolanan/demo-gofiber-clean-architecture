package config

import (
	recover_middleware "github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRecovery() recover_middleware.Config {
	return recover_middleware.Config{
		EnableStackTrace: true,
	}
}
