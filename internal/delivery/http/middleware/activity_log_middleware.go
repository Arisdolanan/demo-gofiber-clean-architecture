package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

func ActivityLogger(activityLogUC usecase.ActivityLogUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process the request first
		err := c.Next()

		// Skip logging if request failed (optional, but usually we log attempts)
		// Or skip if it's the activity log endpoint itself to avoid noise
		path := c.Path()
		if strings.Contains(path, "/activity-logs") {
			return err
		}

		// Only log if user is authenticated
		userIDVal := c.Locals("user_id")
		if userIDVal == nil {
			return err
		}

		userID := userIDVal.(int64)

		// Extract school_id if available
		var schoolID *int64
		if sID := c.Locals("school_id"); sID != nil {
			val := sID.(int64)
			schoolID = &val
		}

		// Identify Module
		module := "General"
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 3 { // e.g., api/v1/users -> users
			module = parts[2]
		}

		// Identify Action
		action := "READ"
		switch c.Method() {
		case fiber.MethodPost:
			action = "CREATE"
		case fiber.MethodPut, fiber.MethodPatch:
			action = "UPDATE"
		case fiber.MethodDelete:
			action = "DELETE"
		}

		// Create Description
		description := fmt.Sprintf("%s request to %s", c.Method(), path)
		if action == "READ" && c.Method() == fiber.MethodGet {
			// Don't log every GET if you want to save space, but the user asked for it
			// description = fmt.Sprintf("Viewed %s", module)
		}

		// Log asynchronously to not block the response
		go func(l *entity.ActivityLog) {
			err := activityLogUC.LogActivity(context.Background(), l)
			if err != nil {
				// This will help us see WHY it fails in the terminal
				fmt.Printf("[ActivityLogger] Failed to log activity: %v\n", err)
			}
		}(&entity.ActivityLog{
			UserID:      &userID,
			SchoolID:    schoolID,
			Action:      action,
			Module:      strings.Title(module),
			Description: description,
			IPAddress:   func() *string { s := c.IP(); return &s }(),
			UserAgent:   func() *string { s := string(c.Request().Header.UserAgent()); return &s }(),
		})

		return err
	}
}
