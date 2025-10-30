package pdfs

import (
	"testing"
)

func TestTemplateManager_RenderTemplate(t *testing.T) {
	templateManager := NewTemplateManager("./")

	invoiceData := TemplateData{
		"InvoiceNumber":   "INV-2023-001",
		"Date":            "2023-10-15",
		"CustomerName":    "John Doe",
		"CustomerAddress": "123 Main St, City, Country",
		"Items": []map[string]interface{}{
			{
				"Description": "Product 1",
				"Quantity":    2,
				"UnitPrice":   "$10.00",
				"Total":       "$20.00",
			},
			{
				"Description": "Product 2",
				"Quantity":    1,
				"UnitPrice":   "$15.00",
				"Total":       "$15.00",
			},
		},
		"TotalAmount": "$35.00",
	}

	htmlContent, err := templateManager.RenderTemplate("invoice", invoiceData)
	if err != nil {
		t.Errorf("Failed to render invoice template: %v", err)
	}

	if htmlContent == "" {
		t.Error("Rendered HTML content is empty")
	}

	if !contains(htmlContent, "INV-2023-001") {
		t.Error("Rendered content does not contain invoice number")
	}

	if !contains(htmlContent, "John Doe") {
		t.Error("Rendered content does not contain customer name")
	}

	if !contains(htmlContent, "Product 1") {
		t.Error("Rendered content does not contain item description")
	}
}

func TestTemplateManager_RenderReceiptTemplate(t *testing.T) {
	templateManager := NewTemplateManager("./")

	receiptData := TemplateData{
		"ReceiptNumber":   "RCT-2023-001",
		"Date":            "2023-10-15",
		"CustomerName":    "Jane Smith",
		"CustomerAddress": "456 Oak Ave, Town, Country",
		"Items": []map[string]interface{}{
			{
				"Name":     "Service A",
				"Quantity": 1,
				"Price":    "$50.00",
			},
		},
		"Total":       "$50.00",
		"CompanyInfo": "ACME Corp - contact@acme.com",
	}

	htmlContent, err := templateManager.RenderTemplate("receipt", receiptData)
	if err != nil {
		t.Errorf("Failed to render receipt template: %v", err)
	}

	if htmlContent == "" {
		t.Error("Rendered HTML content is empty")
	}

	if !contains(htmlContent, "RCT-2023-001") {
		t.Error("Rendered content does not contain receipt number")
	}

	if !contains(htmlContent, "Jane Smith") {
		t.Error("Rendered content does not contain customer name")
	}

	if !contains(htmlContent, "Service A") {
		t.Error("Rendered content does not contain item name")
	}
}

func TestTemplateManager_RenderReportTemplate(t *testing.T) {
	templateManager := NewTemplateManager("./")

	reportData := TemplateData{
		"ReportTitle":   "Sales Report Q4 2023",
		"Period":        "October - December 2023",
		"GeneratedDate": "2023-12-31",
		"Summary":       "This report summarizes sales performance for Q4 2023.",
		"Headers":       []string{"Product", "Units Sold", "Revenue"},
		"DataRows": [][]string{
			{"Product A", "100", "$1,000"},
			{"Product B", "200", "$2,000"},
		},
		"ChartPlaceholder": "Sales chart visualization would appear here",
	}

	htmlContent, err := templateManager.RenderTemplate("report", reportData)
	if err != nil {
		t.Errorf("Failed to render report template: %v", err)
	}

	if htmlContent == "" {
		t.Error("Rendered HTML content is empty")
	}

	if !contains(htmlContent, "Sales Report Q4 2023") {
		t.Error("Rendered content does not contain report title")
	}

	if !contains(htmlContent, "Product A") {
		t.Error("Rendered content does not contain data row")
	}

	if !contains(htmlContent, "Units Sold") {
		t.Error("Rendered content does not contain table headers")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
