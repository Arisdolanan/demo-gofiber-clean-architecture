package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/template/pdfs"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PDFUsecase interface {
	GeneratePDF(req entity.PDFRequest) (*entity.PDFResponse, error)
	GeneratePDFFromTemplate(templateName string, templateData pdfs.TemplateData, filename string) (*entity.PDFResponse, error)
}

type pdfUsecase struct {
	templateManager *pdfs.TemplateManager
	log             *logrus.Logger
}

func NewPDFUsecase(log *logrus.Logger) PDFUsecase {
	// Get PDF configuration from Viper
	pdfConfig := configuration.GetPDFConfig()

	// Set wkhtmltopdf binary path
	if pdfConfig.BinaryPath != "" {
		wkhtmltopdf.SetPath(pdfConfig.BinaryPath)
	} else {
		wkhtmltopdf.SetPath("/usr/local/bin/wkhtmltopdf") // Default path
	}

	// Get template directory from configuration or use default
	templatePath := pdfConfig.TemplateDir
	if templatePath == "" {
		templatePath = "internal/template/pdfs"
	}

	// Ensure the template path is valid
	templatePath = resolveTemplatePath(templatePath)

	templateManager := pdfs.NewTemplateManager(templatePath)

	return &pdfUsecase{
		templateManager: templateManager,
		log:             log,
	}
}

// resolveTemplatePath ensures the template path is valid and resolves correctly across environments
func resolveTemplatePath(configuredPath string) string {
	// If it's already an absolute path that exists, use it
	if filepath.IsAbs(configuredPath) {
		if _, err := os.Stat(configuredPath); err == nil {
			return configuredPath
		}
	}

	// Try common relative paths
	possiblePaths := []string{
		configuredPath,
		"internal/template/pdfs",
		"./internal/template/pdfs",
		"../internal/template/pdfs",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Try to determine the correct path based on current working directory
	if wd, err := os.Getwd(); err == nil {
		attempts := []string{
			filepath.Join(wd, "internal/template/pdfs"),
			filepath.Join(wd, "..", "internal/template/pdfs"),
			filepath.Join(wd, "..", "..", "internal/template/pdfs"),
		}

		for _, attempt := range attempts {
			if _, err := os.Stat(attempt); err == nil {
				return attempt
			}
		}
	}

	// Fall back to default
	return "internal/template/pdfs"
}

func (uc *pdfUsecase) GeneratePDF(req entity.PDFRequest) (*entity.PDFResponse, error) {
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF generator: %w", err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	// Set metadata if provided
	if req.Title != "" {
		pdfg.Title.Set(req.Title)
	}

	// Note: Author, Subject, Keywords are not directly supported as global options
	// They would need to be set as PDF metadata using other methods

	// Create input page from HTML content
	pageReader := wkhtmltopdf.NewPageReader(strings.NewReader(req.HTMLContent))

	// Set page options
	pageReader.FooterRight.Set("[page]")
	pageReader.FooterFontSize.Set(10)

	// Add page to generator
	pdfg.AddPage(pageReader)

	// Create private PDF directory if it doesn't exist (changed to private storage for all PDFs)
	pdfDir := "./storage/private/pdfs"
	if err := os.MkdirAll(pdfDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create private PDF directory: %w", err)
	}

	// Generate PDF
	if err := pdfg.Create(); err != nil {
		return nil, fmt.Errorf("failed to create PDF: %w", err)
	}

	// Generate unique filename
	filename := req.Filename
	if filepath.Ext(filename) == "" {
		filename += ".pdf"
	}

	// If filename doesn't have a UUID, add one to make it unique
	if !strings.Contains(filename, uuid.New().String()[:8]) {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		filename = fmt.Sprintf("%s_%s%s", name, uuid.New().String()[:8], ext)
	}

	// Save PDF to private directory
	pdfPath := filepath.Join(pdfDir, filename)
	if err := pdfg.WriteFile(pdfPath); err != nil {
		return nil, fmt.Errorf("failed to save PDF: %w", err)
	}

	// Return response with private path
	response := &entity.PDFResponse{
		Filename: filename,
		URL:      fmt.Sprintf("/storage/private/pdfs/%s", filename),
	}

	return response, nil
}

func (uc *pdfUsecase) GeneratePDFFromTemplate(templateName string, templateData pdfs.TemplateData, filename string) (*entity.PDFResponse, error) {
	// Render template to HTML
	htmlContent, err := uc.templateManager.RenderTemplate(templateName, templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	// Create PDF request with rendered HTML
	pdfReq := entity.PDFRequest{
		HTMLContent: htmlContent,
		Filename:    filename,
	}

	// Generate PDF using existing method (now stores in private directory)
	return uc.GeneratePDF(pdfReq)
}
