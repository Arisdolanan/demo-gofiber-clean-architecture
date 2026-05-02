package route

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/http/controllers"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/http/middleware"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App                   *fiber.App
	AuthMiddleware        fiber.Handler
	AuthController        *controllers.AuthController
	EmailController       *controllers.EmailController
	UserController        *controllers.UserController
	PDFController         *controllers.PDFController
	ExcelController       *controllers.ExcelController
	FileController        *controllers.FileController
	SchoolController      *controllers.SchoolController
	RBACController        *controllers.RBACController
	AcademicController    *controllers.AcademicController
	PeopleController      *controllers.PeopleController
	OperationController   *controllers.OperationController
	SettingController     *controllers.SettingController
	BackupController      *controllers.BackupController
	ActivityLogController *controllers.ActivityLogController
	ActivityLogUsecase    usecase.ActivityLogUsecase
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

	// Diagnostic route to verify backend version
	c.App.Get("/ping-test", func(ctx *fiber.Ctx) error {
		return ctx.SendString("PONG - Backend Updated: " + ctx.IP())
	})
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
	auth.Post("/switch-school", c.AuthMiddleware, c.AuthController.SwitchSchool)

	// Create a protected API group with authentication and activity logging middleware
	api := c.App.Group("/api/v1", c.AuthMiddleware, middleware.ActivityLogger(c.ActivityLogUsecase))

	// Diagnostic route for protected area
	api.Get("/auth-test", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, Authenticated World!")
	})

	// Activity Logs (protected) - Moved up for priority
	activity := api.Group("/activity-logs")
	activity.Get("", c.ActivityLogController.GetActivities)
	activity.Get("/deletions", c.ActivityLogController.GetDeletions)

	// User routes (protected)
	users := api.Group("/users")
	users.Post("", c.UserController.CreateUser)
	users.Get("", c.UserController.GetAllUsers)
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

	// Settings routes (protected)
	settings := api.Group("/settings")
	settings.Get("", c.SettingController.GetSettings)
	settings.Put("", c.SettingController.UpdateSettings)

	// Backup routes (protected)
	backup := api.Group("/backup")
	backup.Get("/list", c.BackupController.GetBackups)
	backup.Post("/manual", c.BackupController.CreateManualBackup)
	backup.Post("/restore", c.BackupController.RestoreBackup)

	// School & SaaS (protected)
	api.Post("/schools", c.SchoolController.RegisterSchool)
	api.Get("/schools", c.SchoolController.GetAllSchools)
	api.Get("/schools/:id", c.SchoolController.GetSchoolByID)
	api.Put("/schools/:id", c.SchoolController.UpdateSchool)
	api.Post("/schools/packages", c.SchoolController.CreatePackage)
	api.Post("/schools/:id/license", c.SchoolController.AssignLicense)

	// RBAC (protected)
	rbac := api.Group("/rbac")
	rbac.Post("/roles", c.RBACController.CreateRole)
	rbac.Get("/roles", c.RBACController.GetRoles)
	rbac.Post("/roles/:id/permissions", c.RBACController.AssignPermission)
	rbac.Get("/roles/:id/permissions", c.RBACController.GetRolePermissions)
	rbac.Get("/permissions", c.RBACController.GetPermissions)
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
	people.Delete("/teachers/:id", c.PeopleController.DeleteTeacher)

	people.Post("/students", c.PeopleController.CreateStudent)
	people.Put("/students/:id", c.PeopleController.UpdateStudent)
	people.Get("/students/list", c.PeopleController.GetAllStudents)
	people.Get("/students/:id", c.PeopleController.GetStudentByID)
	people.Get("/students/users/:user_id", c.PeopleController.GetStudentByUserID)
	people.Delete("/students/:id", c.PeopleController.DeleteStudent)
	people.Get("/students/:student_id/sections", c.PeopleController.GetStudentSections)
	people.Get("/sections/:section_id/students", c.PeopleController.GetStudentsBySection)
	people.Post("/enroll", c.PeopleController.EnrollStudent)

	// Parent routes (Management now unified in Student Form)
	people.Get("/parents/search", c.PeopleController.SearchParents)
	// standalone parent creation/update removed as per user request

	people.Post("/staff", c.PeopleController.CreateStaff)
	people.Put("/staff/:id", c.PeopleController.UpdateStaff)
	people.Get("/staff", c.PeopleController.GetStaff)
	people.Get("/staff/users/:user_id", c.PeopleController.GetStaffByUserID)
	people.Delete("/staff/:id", c.PeopleController.DeleteStaff)

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
	ops.Get("/attendance/section/:section_id/students", c.OperationController.GetSectionStudentsForAttendance)
	ops.Post("/attendance/section", c.OperationController.RecordSectionAttendance)
	ops.Post("/attendance/teacher", c.OperationController.RecordTeacherAttendance)
	ops.Put("/attendance/teacher/:id", c.OperationController.UpdateTeacherAttendance)
	ops.Get("/teachers/:teacher_id/attendance", c.OperationController.GetTeacherAttendance)
	ops.Post("/attendance/staff", c.OperationController.RecordStaffAttendance)
	ops.Put("/attendance/staff/:id", c.OperationController.UpdateStaffAttendance)
	ops.Get("/staff/:employee_id/attendance", c.OperationController.GetStaffAttendance)

	ops.Get("/notifications", c.OperationController.GetNotifications)
	ops.Get("/report-card", c.OperationController.GetReportCard)
	ops.Get("/report-card/section", c.OperationController.GetSectionReportCards)

	ops.Get("/integrations/definitions", c.OperationController.GetIntegrationDefinitions)
}
