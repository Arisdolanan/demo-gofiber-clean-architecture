# API Documentation

This document describes the API endpoints available in the GoFiber Clean Architecture project.

## PDF Generation API

### Endpoints

#### PDF Generation from HTML
```
POST /api/v1/pdf/generate
```

#### PDF Generation from Template
```
POST /api/v1/pdf/generate-template
```

### Description

Generates PDF files from HTML content or predefined templates. This endpoint follows the clean architecture pattern with proper separation of concerns:

**Note**: All generated PDFs are now stored in the private directory (`storage/private/pdfs/`) and are not accessible via HTTP for security reasons.

### Request Body

#### HTML Content Generation
```json
{
  "html_content": "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title>Sales Report</title><style>body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; } .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); } .header { text-align: center; border-bottom: 2px solid #007bff; padding-bottom: 15px; margin-bottom: 25px; } .header h1 { color: #007bff; margin: 0; } .info-section { display: flex; justify-content: space-between; margin-bottom: 25px; } .info-box { flex: 1; } .info-box h3 { color: #007bff; margin-top: 0; } .summary-box { background-color: #e9f7ff; padding: 15px; border-radius: 5px; margin-bottom: 25px; } .summary-box h3 { margin-top: 0; color: #007bff; } table { width: 100%; border-collapse: collapse; margin: 20px 0; } th, td { border: 1px solid #ddd; padding: 12px; text-align: left; } th { background-color: #007bff; color: white; } tr:nth-child(even) { background-color: #f2f2f2; } tr:hover { background-color: #e9f7ff; } .footer { text-align: center; margin-top: 30px; padding-top: 15px; border-top: 1px solid #eee; color: #666; }</style></head><body><div class=\"container\"><div class=\"header\"><h1>Quarterly Sales Report</h1><p>Q3 2023 | Generated on October 15, 2023</p></div><div class=\"info-section\"><div class=\"info-box\"><h3>Company Information</h3><p><strong>Company:</strong> TechSolutions Inc.</p><p><strong>Address:</strong> 123 Business Avenue, Tech City</p><p><strong>Contact:</strong> info@techsolutions.com</p></div><div class=\"info-box\"><h3>Report Period</h3><p><strong>Quarter:</strong> Q3 2023</p><p><strong>Start Date:</strong> July 1, 2023</p><p><strong>End Date:</strong> September 30, 2023</p></div></div><div class=\"summary-box\"><h3>Executive Summary</h3><p>Q3 2023 showed strong performance with total revenue of $1.2M, representing a 15% increase from Q2 2023. Our flagship product line contributed 65% of total revenue, while new market segments showed promising 25% growth.</p></div><table><thead><tr><th>Product Line</th><th>Q2 2023 Revenue</th><th>Q3 2023 Revenue</th><th>Growth</th><th>% of Total</th></tr></thead><tbody><tr><td>Enterprise Solutions</td><td>$450,000</td><td>$585,000</td><td>$135,000</td><td>48.75%</td></tr><tr><td>Small Business Suite</td><td>$280,000</td><td>$320,000</td><td>$40,000</td><td>26.67%</td></tr><tr><td>Mobile Apps</td><td>$190,000</td><td>$235,000</td><td>$45,000</td><td>19.58%</td></tr><tr><td>Consulting Services</td><td>$80,000</td><td>$60,000</td><td>-$20,000</td><td>5.00%</td></tr></tbody><tfoot><tr style=\"font-weight: bold; background-color: #007bff; color: white;\"><td>Total</td><td>$1,000,000</td><td>$1,200,000</td><td>$200,000</td><td>100%</td></tr></tfoot></table><div class=\"footer\"><p>Confidential - For Internal Use Only</p><p>TechSolutions Inc. | www.techsolutions.com</p></div></div></body></html>",
  "filename": "q3_2023_sales_report",
  "title": "Q3 2023 Sales Report",
  "author": "TechSolutions Inc.",
  "subject": "Quarterly Sales Report Q3 2023",
  "keywords": "sales,report,quarterly,2023,Q3"
}
```

#### Template-based Generation
```json
{
  "template_name": "invoice",
  "template_data": {
    "InvoiceNumber": "INV-2023-001",
    "Date": "2023-10-15",
    "CustomerName": "John Doe",
    "CustomerAddress": "123 Main St, City, Country",
    "Items": [
      {
        "Description": "Web Design Services",
        "Quantity": 10,
        "UnitPrice": "$75.00",
        "Total": "$750.00"
      },
      {
        "Description": "Domain Registration",
        "Quantity": 1,
        "UnitPrice": "$15.00",
        "Total": "$15.00"
      },
      {
        "Description": "Hosting (1 year)",
        "Quantity": 1,
        "UnitPrice": "$120.00",
        "Total": "$120.00"
      }
    ],
    "TotalAmount": "$885.00"
  },
  "filename": "invoice_2023_001",
  "title": "Invoice INV-2023-001"
}
```

#### Fields

##### HTML Content Generation
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| html_content | string | Yes | The HTML content to convert to PDF |
| filename | string | Yes | The filename for the generated PDF |
| title | string | No | The title of the PDF document |
| author | string | No | The author of the PDF document |
| subject | string | No | The subject of the PDF document |
| keywords | string | No | Keywords for the PDF document |

##### Template-based Generation
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| template_name | string | Yes | The name of the template to use (invoice, receipt, report) |
| template_data | object | Yes | The data to populate the template with |
| filename | string | Yes | The filename for the generated PDF |
| title | string | No | The title of the PDF document |

### Response

```json
{
  "status": 200,
  "message": "PDF generated successfully",
  "data": {
    "filename": "test_a1b2c3d4.pdf",
    "url": "/storage/private/pdfs/test_a1b2c3d4.pdf"
  }
}
```

### Authentication

All endpoints require authentication. A valid JWT token must be provided in the Authorization header:

```
Authorization: Bearer <token>
```

### Storage Options

#### Private Storage (All Endpoints)
- All PDFs are stored in `./storage/private/pdfs/` directory
- PDFs are NOT publicly accessible via HTTP
- Suitable for all documents as they are treated as sensitive

**Note**: The previous distinction between public and private storage has been removed. All PDFs are now stored in the private directory for security reasons.

### Example Usage

#### Generate PDF from HTML
```bash
curl -X POST "http://localhost:3000/api/v1/pdf/generate" \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "html_content": "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title>Sales Report</title><style>body { font-family: \"Segoe UI\", Tahoma, Geneva, Verdana, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; } .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); } .header { text-align: center; border-bottom: 2px solid #007bff; padding-bottom: 15px; margin-bottom: 25px; } .header h1 { color: #007bff; margin: 0; } .info-section { display: flex; justify-content: space-between; margin-bottom: 25px; } .info-box { flex: 1; } .info-box h3 { color: #007bff; margin-top: 0; } .summary-box { background-color: #e9f7ff; padding: 15px; border-radius: 5px; margin-bottom: 25px; } .summary-box h3 { margin-top: 0; color: #007bff; } table { width: 100%; border-collapse: collapse; margin: 20px 0; } th, td { border: 1px solid #ddd; padding: 12px; text-align: left; } th { background-color: #007bff; color: white; } tr:nth-child(even) { background-color: #f2f2f2; } tr:hover { background-color: #e9f7ff; } .footer { text-align: center; margin-top: 30px; padding-top: 15px; border-top: 1px solid #eee; color: #666; }</style></head><body><div class=\"container\"><div class=\"header\"><h1>Quarterly Sales Report</h1><p>Q3 2023 | Generated on October 15, 2023</p></div><div class=\"info-section\"><div class=\"info-box\"><h3>Company Information</h3><p><strong>Company:</strong> TechSolutions Inc.</p><p><strong>Address:</strong> 123 Business Avenue, Tech City</p><p><strong>Contact:</strong> info@techsolutions.com</p></div><div class=\"info-box\"><h3>Report Period</h3><p><strong>Quarter:</strong> Q3 2023</p><p><strong>Start Date:</strong> July 1, 2023</p><p><strong>End Date:</strong> September 30, 2023</p></div></div><div class=\"summary-box\"><h3>Executive Summary</h3><p>Q3 2023 showed strong performance with total revenue of $1.2M, representing a 15% increase from Q2 2023. Our flagship product line contributed 65% of total revenue, while new market segments showed promising 25% growth.</p></div><table><thead><tr><th>Product Line</th><th>Q2 2023 Revenue</th><th>Q3 2023 Revenue</th><th>Growth</th><th>% of Total</th></tr></thead><tbody><tr><td>Enterprise Solutions</td><td>$450,000</td><td>$585,000</td><td>$135,000</td><td>48.75%</td></tr><tr><td>Small Business Suite</td><td>$280,000</td><td>$320,000</td><td>$40,000</td><td>26.67%</td></tr><tr><td>Mobile Apps</td><td>$190,000</td><td>$235,000</td><td>$45,000</td><td>19.58%</td></tr><tr><td>Consulting Services</td><td>$80,000</td><td>$60,000</td><td>-$20,000</td><td>5.00%</td></tr></tbody><tfoot><tr style=\"font-weight: bold; background-color: #007bff; color: white;\"><td>Total</td><td>$1,000,000</td><td>$1,200,000</td><td>$200,000</td><td>100%</td></tr></tfoot></table><div class=\"footer\"><p>Confidential - For Internal Use Only</p><p>TechSolutions Inc. | www.techsolutions.com</p></div></div></body></html>",
    "filename": "q3_2023_sales_report.pdf",
    "title": "Q3 2023 Sales Report",
    "author": "TechSolutions Inc.",
    "subject": "Quarterly Sales Report Q3 2023",
    "keywords": "sales,report,quarterly,2023,Q3"
  }'
```

#### Generate PDF from Template
```bash
curl -X POST "http://localhost:3000/api/v1/pdf/generate-template" \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "template_name": "invoice",
    "template_data": {
      "InvoiceNumber": "INV-2023-001",
      "Date": "2023-10-15",
      "CustomerName": "John Doe",
      "CustomerAddress": "123 Main St, City, Country",
      "Items": [
        {
          "Description": "Web Design Services",
          "Quantity": 10,
          "UnitPrice": "$75.00",
          "Total": "$750.00"
        },
        {
          "Description": "Domain Registration",
          "Quantity": 1,
          "UnitPrice": "$15.00",
          "Total": "$15.00"
        },
        {
          "Description": "Hosting (1 year)",
          "Quantity": 1,
          "UnitPrice": "$120.00",
          "Total": "$120.00"
        }
      ],
      "TotalAmount": "$885.00"
    },
    "filename": "invoice_2023_001.pdf"
  }'
```

### Error Handling

The endpoints handle various error conditions:

- **400 Bad Request**: Invalid request body or missing required fields
- **401 Unauthorized**: Missing or invalid authentication token
- **500 Internal Server Error**: PDF generation failed or template rendering failed

### Dependencies

This feature requires the `wkhtmltopdf` binary to be installed on the system. Installation instructions:

#### macOS
```bash
brew install wkhtmltopdf
```

or 

```bash
curl -L https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-2/wkhtmltox-0.12.6-2.macos-cocoa.pkg -O

sudo open wkhtmltox-0.12.6-2.macos-cocoa.pkg

```

#### Ubuntu/Debian
```bash
sudo apt-get install wkhtmltopdf
```

#### CentOS/RHEL
```bash
sudo yum install wkhtmltopdf
```

## Excel Import API

### Endpoint

#### Excel Import
```
POST /api/v1/excel/import
```

### Description

Imports data from Excel files into the database. This endpoint supports both .xlsx and .xls file formats.

### Request Body

The request should be a multipart form with the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| file | file | Yes | The Excel file to import |
| sheet | string | No | The name of the Excel sheet to import (defaults to first sheet) |
| table | string | Yes | The name of the database table to import data into |

### Response

```json
{
  "status": 200,
  "message": "Excel data imported successfully",
  "data": {
    "headers": ["id", "name", "email"],
    "rows": [
      {"id": "1", "name": "John Doe", "email": "john@example.com"},
      {"id": "2", "name": "Jane Smith", "email": "jane@example.com"}
    ]
  }
}
```

### Authentication

All endpoints require authentication. A valid JWT token must be provided in the Authorization header:

```
Authorization: Bearer <token>
```

### Example Usage

```bash
curl -X POST "http://localhost:3000/api/v1/excel/import" \
  -H "Authorization: Bearer your-jwt-token" \
  -F "file=@users.xlsx" \
  -F "sheet=Sheet1" \
  -F "table=users"
```

### Error Handling

The endpoint handles various error conditions:

- **400 Bad Request**: Invalid request body, missing required fields, or invalid file format
- **401 Unauthorized**: Missing or invalid authentication token
- **500 Internal Server Error**: Excel import failed or file processing error

### Security

For security reasons, only Excel import functionality is available through the API. Excel export endpoints have been removed to prevent unauthorized data extraction.

## Using Excel Functionality in Other Usecases

While Excel export endpoints have been removed from the public API for security reasons, other usecases can still utilize the Excel export functionality directly through the ExcelUsecase interface. This approach provides better control and security as it doesn't expose data export capabilities through HTTP endpoints.

### Available Methods

The ExcelUsecase interface provides the following methods for exporting data to Excel:

1. `ExportExcel(req entity.ExcelExportRequest) (*entity.ExcelResponse, error)`
   - Exports data to Excel file format
   - Accepts various data types (structs, slices, maps)

2. `ExportCustomQueryToExcel(db *sqlx.DB, req entity.CustomQueryRequest) (*entity.ExcelResponse, error)`
   - Executes a custom SQL query and exports results to Excel
   - Only allows SELECT statements for security
   - Uses the centralized database connection

3. `ExportQueryResultToExcel(queryResult *entity.QueryResult, sheet, filename string) (*entity.ExcelResponse, error)`
   - Exports pre-executed query results to Excel
   - Useful when another usecase has already executed a query

### Usage Examples

#### 1. Direct Data Export
```go
// In another usecase, inject ExcelUsecase
type SomeOtherUseCase struct {
    excelUsecase usecase.ExcelUsecase
}

func (s *SomeOtherUseCase) ExportUserData() error {
    // Prepare data to export
    users := []map[string]interface{}{
        {"id": 1, "name": "John Doe", "email": "john@example.com"},
        {"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
    }

    // Create export request
    req := entity.ExcelExportRequest{
        Data:     users,
        Sheet:    "Users",
        Filename: "users_report.xlsx",
        Table:    "users",
    }

    // Export to Excel
    result, err := s.excelUsecase.ExportExcel(req)
    if err != nil {
        return fmt.Errorf("failed to export users to Excel: %w", err)
    }

    // Handle the result (e.g., store file path, send notification, etc.)
    fmt.Printf("Excel file created: %s\n", result.URL)
    return nil
}
```

#### 2. Custom Query Export
```go
// In another usecase with database access
type ReportUseCase struct {
    excelUsecase usecase.ExcelUsecase
    db          *sqlx.DB
}

func (r *ReportUseCase) GenerateSalesReport() error {
    // Create custom query request
    req := entity.CustomQueryRequest{
        Query:    "SELECT u.id, u.username, p.title as profile_title FROM users u LEFT JOIN profiles p ON u.id = p.user_id WHERE u.created_at > $1",
        Params:   map[string]interface{}{"created_at": "2023-01-01"},
        Sheet:    "UsersWithProfiles",
        Filename: "users_with_profiles_report.xlsx",
    }

    // Export query results to Excel
    result, err := r.excelUsecase.ExportCustomQueryToExcel(r.db, req)
    if err != nil {
        return fmt.Errorf("failed to export sales report: %w", err)
    }

    // Handle the result
    fmt.Printf("Sales report created: %s\n", result.URL)
    return nil
}
```

#### 3. Pre-executed Query Result Export
```go
// In another usecase that has already processed data
type AnalyticsUseCase struct {
    excelUsecase usecase.ExcelUsecase
}

func (a *AnalyticsUseCase) ExportProcessedData() error {
    // Assume we have processed data from some analytics operation
    queryResult := &entity.QueryResult{
        Columns: []string{"metric", "value", "date"},
        Rows: []map[string]interface{}{
            {"metric": "page_views", "value": 12500, "date": "2023-10-01"},
            {"metric": "unique_visitors", "value": 8500, "date": "2023-10-01"},
            {"metric": "conversion_rate", "value": 0.035, "date": "2023-10-01"},
        },
    }

    // Export the processed data directly
    result, err := a.excelUsecase.ExportQueryResultToExcel(queryResult, "Analytics", "analytics_report.xlsx")
    if err != nil {
        return fmt.Errorf("failed to export analytics data: %w", err)
    }

    // Handle the result
    fmt.Printf("Analytics report created: %s\n", result.URL)
    return nil
}
```

### Benefits of This Approach

1. **Security**: Data export functionality is not exposed through public API endpoints
2. **Centralized Database Connection**: All database operations use the centralized connection from main.go
3. **Reusability**: Excel export functionality can be used by any usecase in the application
4. **Flexibility**: Supports various data formats and custom SQL queries
5. **Control**: Each usecase can control when and how data is exported
6. **Testability**: Easy to test with mocking
7. **Maintainability**: Changes to Excel export logic only need to be made in one place