package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/sirupsen/logrus"
)

// EmailService interface for sending emails
type EmailService interface {
	SendVerificationEmail(userEmail, token string) error
	SendPasswordResetEmail(userEmail, token string) error
	SendEmail(emailTemplate *entity.EmailTemplate) error
}

// emailService implements EmailService
type emailService struct {
	config configuration.EmailConfig
	logger *logrus.Logger
}

// NewEmailService creates a new email service
func NewEmailService(logger *logrus.Logger) EmailService {
	config := configuration.GetEmailConfig()
	return &emailService{
		config: config,
		logger: logger,
	}
}

// SendVerificationEmail sends email verification email
func (e *emailService) SendVerificationEmail(userEmail, token string) error {
	baseURL := e.config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:3000" // Fallback
	}
	verificationURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", baseURL, token)

	data := entity.EmailVerificationData{
		UserEmail:       userEmail,
		VerificationURL: verificationURL,
		AppName:         configuration.GetAppName(),
		ExpirationHours: e.config.VerificationTokenExpiry / 3600, // Convert seconds to hours
	}

	emailTemplate := &entity.EmailTemplate{
		To:       userEmail,
		Subject:  fmt.Sprintf("Verify your email for %s", data.AppName),
		Template: "verification.html",
		Data: map[string]interface{}{
			"UserEmail":       data.UserEmail,
			"VerificationURL": data.VerificationURL,
			"AppName":         data.AppName,
			"ExpirationHours": data.ExpirationHours,
		},
	}

	return e.SendEmail(emailTemplate)
}

// SendPasswordResetEmail sends password reset email
func (e *emailService) SendPasswordResetEmail(userEmail, token string) error {
	baseURL := e.config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:3000" // Fallback
	}
	resetURL := fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", baseURL, token)

	data := entity.PasswordResetData{
		UserEmail:      userEmail,
		ResetURL:       resetURL,
		AppName:        configuration.GetAppName(),
		ExpirationMins: e.config.ResetTokenExpiry / 60, // Convert seconds to minutes
	}

	emailTemplate := &entity.EmailTemplate{
		To:       userEmail,
		Subject:  fmt.Sprintf("Reset your password for %s", data.AppName),
		Template: "password_reset.html",
		Data: map[string]interface{}{
			"UserEmail":      data.UserEmail,
			"ResetURL":       data.ResetURL,
			"AppName":        data.AppName,
			"ExpirationMins": data.ExpirationMins,
		},
	}

	return e.SendEmail(emailTemplate)
}

// SendEmail sends an email using the provided template
func (e *emailService) SendEmail(emailTemplate *entity.EmailTemplate) error {
	tmplPath := filepath.Join(e.config.TemplateDir, emailTemplate.Template)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		e.logger.Errorf("Error parsing email template %s: %v", tmplPath, err)
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, emailTemplate.Data); err != nil {
		e.logger.Errorf("Error executing email template: %v", err)
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	message := fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		e.config.FromName,
		e.config.FromEmail,
		emailTemplate.To,
		emailTemplate.Subject,
		body.String(),
	)

	// Send email asynchronously in a goroutine
	go func() {
		auth := smtp.PlainAuth("", e.config.SMTPUsername, e.config.SMTPPassword, e.config.SMTPHost)

		addr := fmt.Sprintf("%s:%d", e.config.SMTPHost, e.config.SMTPPort)
		err := smtp.SendMail(addr, auth, e.config.FromEmail, []string{emailTemplate.To}, []byte(message))
		if err != nil {
			e.logger.Errorf("Error sending email to %s: %v", emailTemplate.To, err)
			return
		}

		e.logger.Infof("Email sent successfully to %s", emailTemplate.To)
	}()

	return nil
}
