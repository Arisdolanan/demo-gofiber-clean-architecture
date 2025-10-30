package route

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/http/controllers"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App             *fiber.App
	AuthMiddleware  fiber.Handler
	AuthController  *controllers.AuthController
	EmailController *controllers.EmailController
	UserController  *controllers.UserController
	PDFController   *controllers.PDFController
	ExcelController *controllers.ExcelController
	FileController  *controllers.FileController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(response.HTTPSuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Welcome to GoFiber Clean Architecture",
			Data:    nil,
		})
	})

	// Metrics and health check routes
	c.App.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	// Swagger docs
	c.App.Get("/swagger/*", swagger.HandlerDefault) // default

	// Static files
	c.App.Static("/storage", "./storage")

	// Public file routes (no authentication required)
	c.App.Get("/api/v1/files/public", c.FileController.GetPublicFiles)
}

func (c *RouteConfig) SetupAuthRoute() {
	// Guest auth routes (no middleware required)
	auth := c.App.Group("/api/v1/auth")
	auth.Post("/register", c.AuthController.Register)
	auth.Post("/login", c.AuthController.Login)
	auth.Post("/refresh", c.AuthController.Refresh)

	// Email verification and password reset routes (no authentication required)
	auth.Get("/verify-email", c.EmailController.VerifyEmail)
	auth.Post("/resend-verification", c.EmailController.ResendVerificationEmail)
	auth.Post("/forgot-password", c.EmailController.RequestPasswordReset)
	auth.Post("/reset-password", c.EmailController.ResetPassword)

	// Protected auth routes (with JWT middleware)
	auth.Get("/verify", c.AuthMiddleware, c.AuthController.Verify)
	auth.Post("/logout", c.AuthMiddleware, c.AuthController.Logout)

	// Create a protected API group with authentication middleware
	api := c.App.Group("/api/v1", c.AuthMiddleware)

	// Test route for protected area
	api.Get("/auth", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	// User routes (protected)
	users := api.Group("/users")
	users.Post("/", c.UserController.CreateUser)
	users.Get("/", c.UserController.GetAllUsers)
	users.Get("/:id", c.UserController.GetUserByID)
	users.Put("/:id", c.UserController.UpdateUser)
	users.Delete("/:id", c.UserController.DeleteUser)
	users.Delete("/:id/soft-delete", c.UserController.SoftDeleteUser)

	// PDF generation routes (protected)
	pdf := api.Group("/pdf")
	pdf.Post("/generate", c.PDFController.GeneratePDF)
	pdf.Post("/generate-template", c.PDFController.GeneratePDFFromTemplate)

	// Excel routes (protected)
	excel := api.Group("/excel")
	excel.Post("/import", c.ExcelController.ImportExcel)
	excel.Post("/export", c.ExcelController.ExportExcel)

	// File management routes (protected)
	files := api.Group("/files")
	files.Post("/upload", c.FileController.UploadFile)
	files.Get("/", c.FileController.GetUserFiles)
	files.Get("/private", c.FileController.GetPrivateFiles)
	files.Put("/:id", c.FileController.UpdateFile)
	files.Delete("/:id", c.FileController.DeleteFile)
	files.Get("/:id/download", c.FileController.DownloadFile)

	// Admin routes (protected)
	admin := api.Group("/admin")
	admin.Post("/cleanup-tokens", c.EmailController.CleanupExpiredTokens)
}
