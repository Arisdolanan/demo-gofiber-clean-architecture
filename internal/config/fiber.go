package config

import (
	"log"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func NewFiber() *fiber.App {
	appConfig := configuration.GetAppConfig()

	app := fiber.New(fiber.Config{
		AppName:      appConfig.Name,
		ServerHeader: "Fiber",
		ErrorHandler: NewErrorHandler(),
		Prefork:      appConfig.Prefork,
	})
	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		errorLogConfig := utils.LogDev()
		logger := log.New(errorLogConfig.Output, "", log.LstdFlags)
		logger.Printf("[ERROR]: %s - Path: %s - Method: %s - Code: %d\n", err.Error(), ctx.Path(), ctx.Method(), code)

		return ctx.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
