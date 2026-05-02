package config

import (
	"github.com/IBM/sarama"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/http/controllers"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/http/middleware"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/http/route"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/messaging/kafka"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/cache"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/email"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/messaging/rabbitmq"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/redis"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover_middleware "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	DB               *sqlx.DB
	Redis            *cache.RedisCache
	App              *fiber.App
	Log              *logrus.Logger
	Validate         *validator.Validate
	KafkaProducer    sarama.SyncProducer
	RabbitMQProducer *rabbitmq.Producer
}

func Bootstrap(cfg *BootstrapConfig) {
	loggingConfig := configuration.GetLoggingConfig()

	// Setup logging middleware based on configuration from Viper
	if loggingConfig.UseLogrus {
		cfg.App.Use(logger.New(utils.LogDev()))
	} else {
		cfg.App.Use(logger.New(utils.LogFiber()))
	}

	// Setup other middleware
	cfg.App.Use(cors.New(SetupCors()))
	cfg.App.Use(limiter.New(SetupLimiter()))
	cfg.App.Use(recover_middleware.New(SetupRecovery()))

	// Initialize repositories
	authRepo := postgresql.NewAuthRepository(cfg.DB)
	userRepo := postgresql.NewUserRepository(cfg.DB)
	fileRepo := postgresql.NewFileRepository(cfg.DB)
	emailRepo := postgresql.NewEmailRepository(cfg.DB)
	authRedisRepo := redis.NewAuthRedisRepository(cfg.Redis.GetClient())
	schoolRepo := postgresql.NewSchoolRepository(cfg.DB)
	roleRepo := postgresql.NewRoleRepository(cfg.DB)
	permRepo := postgresql.NewPermissionRepository(cfg.DB)
	academicRepo := postgresql.NewAcademicRepository(cfg.DB)
	peopleRepo := postgresql.NewPeopleRepository(cfg.DB)
	operationRepo := postgresql.NewOperationRepository(cfg.DB)
	settingRepo := postgresql.NewSettingRepository(cfg.DB)
	backupRepo := postgresql.NewBackupRepository(cfg.DB.DB)
	activityLogRepo := postgresql.NewActivityLogRepository(cfg.DB)

	// Get JWT secret from Viper
	jwtSecret := configuration.GetJWTSecret()

	// Initialize Kafka user producer if Kafka producer is available
	var kafkaUserProducer *kafka.UserProducer
	if cfg.KafkaProducer != nil {
		kafkaUserProducer = kafka.NewUserProducer(cfg.KafkaProducer, cfg.Log)
	}

	// Initialize email service
	emailService := email.NewEmailService(cfg.Log)

	// Initialize usecase
	emailUsecase := usecase.NewEmailUsecase(emailRepo, userRepo, emailService, cfg.Log)
	activityLogUsecase := usecase.NewActivityLogUsecase(activityLogRepo, cfg.Log)
	authUsecase := usecase.NewAuthUsecase(authRepo, authRedisRepo, emailUsecase, cfg.Validate, cfg.Log, jwtSecret, kafkaUserProducer, activityLogUsecase)
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, cfg.Redis, cfg.Log, cfg.Validate)
	fileUsecase := usecase.NewFileUseCase(fileRepo, cfg.Log, cfg.Validate)
	pdfUsecase := usecase.NewPDFUsecase(cfg.Log)
	excelUsecase := usecase.NewExcelUsecase(cfg.Log)
	schoolUsecase := usecase.NewSchoolUsecase(schoolRepo, cfg.Validate, cfg.Log)
	rbacUsecase := usecase.NewRBACUsecase(roleRepo, permRepo, cfg.Validate, cfg.Log)
	academicUsecase := usecase.NewAcademicUsecase(academicRepo, cfg.Validate, cfg.Log)
	peopleUsecase := usecase.NewPeopleUsecase(peopleRepo, userUsecase, cfg.Validate, cfg.Log)
	operationUsecase := usecase.NewOperationUsecase(operationRepo, peopleRepo, cfg.Validate, cfg.Log)
	settingUsecase := usecase.NewSettingUseCase(settingRepo, cfg.Log, cfg.Validate)
	backupUsecase := usecase.NewBackupUseCase(backupRepo, cfg.Log)

	// Initialize controller
	authController := controllers.NewAuthController(authUsecase, cfg.Validate, cfg.Log)
	emailController := controllers.NewEmailController(emailUsecase, cfg.Validate, cfg.Log)
	userController := controllers.NewUserController(userUsecase, cfg.Log)
	fileController := controllers.NewFileController(fileUsecase, cfg.Validate, cfg.Log)
	pdfController := controllers.NewPDFController(pdfUsecase, cfg.Validate, cfg.Log)
	excelController := controllers.NewExcelController(excelUsecase, cfg.Validate, cfg.Log, cfg.DB)
	schoolController := controllers.NewSchoolController(schoolUsecase, cfg.Log)
	rbacController := controllers.NewRBACController(rbacUsecase, cfg.Log)
	academicController := controllers.NewAcademicController(academicUsecase, cfg.Log)
	peopleController := controllers.NewPeopleController(peopleUsecase, cfg.Log)
	operationController := controllers.NewOperationController(operationUsecase, cfg.Log)
	settingController := controllers.NewSettingController(settingUsecase, cfg.Log)
	backupController := controllers.NewBackupController(backupUsecase, cfg.Log)
	activityLogController := controllers.NewActivityLogController(activityLogUsecase, cfg.Log)

	// Initialize middleware with blacklisting support
	authMiddleware := middleware.JWTProtectedWithBlacklist(jwtSecret, authRedisRepo)
	// Alternative: Use basic middleware without blacklisting
	// authMiddleware := middleware.JWTProtected(jwtSecret)

	// Setup routes
	routeConfig := &route.RouteConfig{
		App:                 cfg.App,
		AuthMiddleware:      authMiddleware,
		AuthController:      authController,
		EmailController:     emailController,
		UserController:      userController,
		FileController:      fileController,
		PDFController:       pdfController,
		ExcelController:     excelController,
		SchoolController:    schoolController,
		RBACController:      rbacController,
		AcademicController:  academicController,
		PeopleController:    peopleController,
		OperationController: operationController,
		SettingController:   settingController,
		BackupController:    backupController,
		ActivityLogController: activityLogController,
		ActivityLogUsecase:    activityLogUsecase,
	}
	routeConfig.Setup()
}
