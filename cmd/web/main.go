package main

import (
	_ "github.com/arisdolanan/demo-gofiber-clean-architecture/api/docs"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/config"
	_ "github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/http/controllers"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/cache"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/database"
	kafkaInfra "github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/messaging/kafka"
	rabbitInfra "github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/messaging/rabbitmq"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	_ "github.com/swaggo/fiber-swagger"
)

// @title Fiber Swagger Example API
// @version 1.0
// @description This is a sample server for using Swagger with Fiber.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

func main() {
	viperConfig, err := configuration.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}
	logConfig := utils.CreateLogConfigFromViper()
	log := config.SetupLogRus(logConfig)

	// Initialize Messaging producer
	producer := kafkaInfra.NewKafkaProducer(viperConfig, log)
	rabbitProducer := rabbitInfra.NewRabbitMqProducer(log)

	app := config.NewFiber()
	validate := config.NewValidator()

	config.Bootstrap(&config.BootstrapConfig{
		App:              app,
		Validate:         validate,
		Log:              log,
		DB:               database.ConnectPostgresqlx(),
		Redis:            cache.NewRedisCache(),
		KafkaProducer:    producer,
		RabbitMQProducer: rabbitProducer,
	})

	app.Listen(":" + configuration.GetAppPort())
}
