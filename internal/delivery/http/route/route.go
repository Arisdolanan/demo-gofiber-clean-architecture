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
	FileController      *controllers.FileController
	SchoolController    *controllers.SchoolController
	RBACController      *controllers.RBACController
	AcademicController  *controllers.AcademicController
	PeopleController    *controllers.PeopleController
	OperationController *controllers.OperationController
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
	auth.Post("/verify-email", c.EmailController.VerifyEmail)
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

	// School & SaaS (protected)
	schools := api.Group("/schools")
	schools.Post("/", c.SchoolController.RegisterSchool)
	schools.Get("/:id", c.SchoolController.GetSchoolByID)
	schools.Post("/packages", c.SchoolController.CreatePackage)
	schools.Post("/:id/license", c.SchoolController.AssignLicense)

	// RBAC (protected)
	rbac := api.Group("/rbac")
	rbac.Post("/roles", c.RBACController.CreateRole)
	rbac.Get("/roles", c.RBACController.GetRoles)
	rbac.Post("/roles/:id/permissions", c.RBACController.AssignPermission)
	rbac.Post("/users/:user_id/roles", c.RBACController.AssignRoleToUser)
	rbac.Get("/users/:user_id/roles", c.RBACController.GetUserRoles)

	// Academic (protected)
	academic := api.Group("/academic")
	academic.Post("/sessions", c.AcademicController.CreateSession)
	academic.Put("/sessions/:id", c.AcademicController.UpdateSession)
	academic.Delete("/sessions/:id", c.AcademicController.DeleteSession)
	academic.Get("/sessions", c.AcademicController.GetSessions)
	academic.Get("/sessions/active", c.AcademicController.GetActiveSession)

	academic.Post("/classes", c.AcademicController.CreateClass)
	academic.Put("/classes/:id", c.AcademicController.UpdateClass)
	academic.Delete("/classes/:id", c.AcademicController.DeleteClass)
	academic.Get("/classes", c.AcademicController.GetClasses)

	academic.Post("/subjects", c.AcademicController.CreateSubject)
	academic.Put("/subjects/:id", c.AcademicController.UpdateSubject)
	academic.Delete("/subjects/:id", c.AcademicController.DeleteSubject)
	academic.Get("/subjects", c.AcademicController.GetSubjects)

	academic.Post("/sections", c.AcademicController.CreateSection)
	academic.Put("/sections/:id", c.AcademicController.UpdateSection)
	academic.Delete("/sections/:id", c.AcademicController.DeleteSection)
	academic.Get("/sections", c.AcademicController.GetSections)

	// People (protected)
	people := api.Group("/people")
	people.Post("/teachers", c.PeopleController.CreateTeacher)
	people.Put("/teachers/:id", c.PeopleController.UpdateTeacher)
	people.Get("/teachers", c.PeopleController.GetTeachers)
	people.Get("/teachers/users/:user_id", c.PeopleController.GetTeacherByUserID)
	
	people.Post("/students", c.PeopleController.CreateStudent)
	people.Put("/students/:id", c.PeopleController.UpdateStudent)
	people.Get("/students/list", c.PeopleController.GetAllStudents)
	people.Get("/students/users/:user_id", c.PeopleController.GetStudentByUserID)
	people.Get("/students/:student_id/sections", c.PeopleController.GetStudentSections)
	people.Get("/sections/:section_id/students", c.PeopleController.GetStudentsBySection)
	people.Post("/enroll", c.PeopleController.EnrollStudent)
	
	people.Post("/parents", c.PeopleController.CreateParent)
	people.Post("/parents/link", c.PeopleController.LinkParentToStudent)

	// Operations (protected)
	ops := api.Group("/operations")
	ops.Post("/schedules", c.OperationController.CreateSchedule)
	ops.Put("/schedules/:id", c.OperationController.UpdateSchedule)
	ops.Delete("/schedules/:id", c.OperationController.DeleteSchedule)
	ops.Get("/schedules", c.OperationController.GetAllSchedules)

	ops.Post("/exams", c.OperationController.CreateExam)
	ops.Put("/exams/:id", c.OperationController.UpdateExam)
	ops.Delete("/exams/:id", c.OperationController.DeleteExam)
	ops.Get("/exams", c.OperationController.GetAllExams)
	ops.Get("/exams/:exam_id/marks", c.OperationController.GetExamMarks)

	ops.Post("/marks", c.OperationController.UpdateMark)
	ops.Get("/students/:student_id/marks", c.OperationController.GetStudentMarks)

	ops.Post("/attendance/student", c.OperationController.RecordStudentAttendance)
	ops.Put("/attendance/student/:id", c.OperationController.UpdateStudentAttendance)
	ops.Get("/students/:student_id/attendance", c.OperationController.GetStudentAttendance)
	ops.Get("/attendance", c.OperationController.GetAttendanceReport)
	ops.Post("/attendance/teacher", c.OperationController.RecordTeacherAttendance)

	ops.Get("/notifications", c.OperationController.GetNotifications)
	ops.Get("/report-card", c.OperationController.GetReportCard)
	ops.Get("/report-card/section", c.OperationController.GetSectionReportCards)
}
