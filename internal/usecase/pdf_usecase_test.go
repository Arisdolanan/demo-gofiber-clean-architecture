package usecase

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/template/pdfs"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPDFUsecase_GeneratePDF(t *testing.T) {
	// Create a new PDF usecase with logger
	logger := logrus.New()
	pdfUsecase := NewPDFUsecase(logger)

	// Test data
	req := entity.PDFRequest{
		HTMLContent: "<h1>Test PDF</h1><p>This is a test PDF generated from HTML content.</p>",
		Filename:    "test_document_public",
		Title:       "Test Document",
	}

	// Generate PDF
	response, err := pdfUsecase.GeneratePDF(req)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Filename)
	assert.NotEmpty(t, response.URL)

	// Check if file exists in private directory (all PDFs now go to private directory)
	privatePDFPath := filepath.Join("./storage/private/pdfs", response.Filename)
	if _, err := os.Stat(privatePDFPath); os.IsNotExist(err) {
		t.Errorf("PDF file was not created in private directory: %v", err)
	}

	// Clean up test file
	defer func() {
		_ = os.Remove(privatePDFPath)
	}()
}

func TestPDFUsecase_GeneratePDFFromTemplate(t *testing.T) {
	// Create a new PDF usecase using the actual constructor to ensure consistency
	logger := logrus.New()
	pdfUsecase := NewPDFUsecase(logger)

	// Test data
	templateData := pdfs.TemplateData{
		"InvoiceNumber":   "INV-2023-001",
		"Date":            "2023-10-15",
		"CustomerName":    "John Doe",
		"CustomerAddress": "123 Main St, City, Country",
		"Items": []map[string]interface{}{
			{
				"Description": "Web Design Services",
				"Quantity":    10,
				"UnitPrice":   "$75.00",
				"Total":       "$750.00",
			},
			{
				"Description": "Domain Registration",
				"Quantity":    1,
				"UnitPrice":   "$15.00",
				"Total":       "$15.00",
			},
		},
		"TotalAmount": "$765.00",
	}

	// Generate PDF from template
	response, err := pdfUsecase.GeneratePDFFromTemplate("invoice", templateData, "test_invoice")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Filename)
	assert.NotEmpty(t, response.URL)

	// Check if file exists in private directory
	privatePDFPath := filepath.Join("./storage/private/pdfs", response.Filename)
	if _, err := os.Stat(privatePDFPath); os.IsNotExist(err) {
		t.Errorf("PDF file was not created in private directory: %v", err)
	}

	// Clean up test file
	defer func() {
		_ = os.Remove(privatePDFPath)
	}()
}
