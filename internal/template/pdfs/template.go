package pdfs

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
)

// TemplateData represents the data structure for template rendering
type TemplateData map[string]interface{}

// TemplateManager manages PDF templates
type TemplateManager struct {
	templatesDir string
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(templatesDir string) *TemplateManager {
	return &TemplateManager{
		templatesDir: templatesDir,
	}
}

// RenderTemplate renders a template with data and returns HTML content
func (tm *TemplateManager) RenderTemplate(templateName string, data TemplateData) (string, error) {
	// Construct the full path to the template file
	templatePath := filepath.Join(tm.templatesDir, templateName+".html")

	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	// Execute the template with provided data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}
