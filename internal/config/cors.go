package config

import "github.com/gofiber/fiber/v2/middleware/cors"

func SetupCors() cors.Config {
	return cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		ExposeHeaders:    "Content-Length,Content-Type",
	}
}
