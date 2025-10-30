package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

// ExcelUsecase interface defines the methods for Excel operations
type ExcelUsecase interface {
	ImportExcel(req entity.ExcelImportRequest) (*entity.ExcelData, error)
	ExportExcel(req entity.ExcelExportRequest) (*entity.ExcelResponse, error)
	ExportCustomQueryToExcel(db *sqlx.DB, req entity.CustomQueryRequest) (*entity.ExcelResponse, error)
	ExportQueryResultToExcel(queryResult *entity.QueryResult, sheet, filename string) (*entity.ExcelResponse, error) // Fungsi baru
}

type excelUsecase struct {
	log *logrus.Logger
}

// NewExcelUsecase creates a new instance of ExcelUsecase
func NewExcelUsecase(log *logrus.Logger) ExcelUsecase {
	return &excelUsecase{
		log: log,
	}
}

// ImportExcel imports data from an Excel file
func (uc *excelUsecase) ImportExcel(req entity.ExcelImportRequest) (*entity.ExcelData, error) {
	uc.log.Infof("Starting Excel import for table: %s", req.Table)

	// Create a temporary file from the byte data
	tempFile, err := os.CreateTemp("", "excel_import_*.xlsx")
	if err != nil {
		uc.log.Errorf("Failed to create temporary file: %v", err)
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write the file data to the temporary file
	if _, err := tempFile.Write(req.FileData); err != nil {
		return nil, fmt.Errorf("failed to write file data to temporary file: %w", err)
	}

	// Open the Excel file
	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			uc.log.Errorf("failed to close Excel file: %v", err)
		}
	}()

	// Get the sheet name, use the first sheet if not specified
	sheetName := req.Sheet
	if sheetName == "" {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("no sheets found in Excel file")
		}
		sheetName = sheets[0]
	}

	// Read all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from sheet %s: %w", sheetName, err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no data found in sheet %s", sheetName)
	}

	// Extract headers from the first row
	headers := rows[0]

	// Process data rows
	var excelRows []entity.ExcelRow
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		excelRow := make(entity.ExcelRow)

		// Map each cell to its corresponding header
		for j, cellValue := range row {
			if j < len(headers) {
				header := headers[j]
				excelRow[header] = cellValue
			} else {
				// Handle columns without headers
				excelRow[fmt.Sprintf("Column_%d", j+1)] = cellValue
			}
		}

		excelRows = append(excelRows, excelRow)
	}

	uc.log.Infof("Successfully imported Excel data with %d rows", len(excelRows))

	return &entity.ExcelData{
		Headers: headers,
		Rows:    excelRows,
	}, nil
}

// ExportExcel exports data to an Excel file
func (uc *excelUsecase) ExportExcel(req entity.ExcelExportRequest) (*entity.ExcelResponse, error) {
	// Create a new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			uc.log.Errorf("failed to close Excel file: %v", err)
		}
	}()

	uc.log.Infof("Starting Excel export with filename: %s", req.Filename)

	// Get the sheet name, use "Sheet1" if not specified
	sheetName := req.Sheet
	if sheetName == "" {
		sheetName = "Sheet1"
	}

	// Create the sheet if it doesn't exist
	f.NewSheet(sheetName)

	// Convert data to rows based on the data type
	var dataRows []entity.ExcelRow
	var headers []string

	switch data := req.Data.(type) {
	case *entity.ExcelData:
		// If data is already in ExcelData format
		dataRows = data.Rows
		headers = data.Headers
	case []map[string]interface{}:
		// If data is a slice of maps, convert to ExcelData format
		if len(data) > 0 {
			// Extract headers from the first row
			for key := range data[0] {
				headers = append(headers, key)
			}

			// Convert each map to ExcelRow
			for _, row := range data {
				excelRow := make(entity.ExcelRow)
				for key, value := range row {
					excelRow[key] = value
				}
				dataRows = append(dataRows, excelRow)
			}
		}
	case []entity.ExcelRow:
		// If data is already in ExcelRow format
		dataRows = data
		// Extract headers from the first row if available
		if len(data) > 0 {
			for key := range data[0] {
				headers = append(headers, key)
			}
		}
	default:
		// Handle struct data using reflection
		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Slice:
			// Handle slice of structs
			if v.Len() > 0 {
				firstElement := v.Index(0)
				if firstElement.Kind() == reflect.Ptr {
					firstElement = firstElement.Elem()
				}

				if firstElement.Kind() == reflect.Struct {
					// Extract field names as headers
					t := firstElement.Type()
					for i := 0; i < t.NumField(); i++ {
						field := t.Field(i)
						jsonTag := field.Tag.Get("json")
						if jsonTag != "" && jsonTag != "-" {
							headers = append(headers, strings.Split(jsonTag, ",")[0])
						} else {
							headers = append(headers, field.Name)
						}
					}

					// Convert each struct to ExcelRow
					for i := 0; i < v.Len(); i++ {
						elem := v.Index(i)
						if elem.Kind() == reflect.Ptr {
							elem = elem.Elem()
						}

						if elem.Kind() == reflect.Struct {
							excelRow := make(entity.ExcelRow)
							for j, header := range headers {
								field := elem.Field(j)
								if field.IsValid() && field.CanInterface() {
									excelRow[header] = field.Interface()
								} else {
									excelRow[header] = ""
								}
							}
							dataRows = append(dataRows, excelRow)
						}
					}
				}
			}
		case reflect.Struct:
			// Handle single struct
			t := v.Type()
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				jsonTag := field.Tag.Get("json")
				if jsonTag != "" && jsonTag != "-" {
					headers = append(headers, strings.Split(jsonTag, ",")[0])
				} else {
					headers = append(headers, field.Name)
				}
			}

			excelRow := make(entity.ExcelRow)
			for j, header := range headers {
				field := v.Field(j)
				if field.IsValid() && field.CanInterface() {
					excelRow[header] = field.Interface()
				} else {
					excelRow[header] = ""
				}
			}
			dataRows = append(dataRows, excelRow)
		}
	}

	// Write headers to the first row
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to generate cell name for header: %w", err)
		}
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data rows
	for rowIndex, row := range dataRows {
		for colIndex, header := range headers {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2) // +2 because of header row
			if err != nil {
				return nil, fmt.Errorf("failed to generate cell name for data: %w", err)
			}

			if value, exists := row[header]; exists {
				f.SetCellValue(sheetName, cell, value)
			}
		}
	}

	// Ensure storage directory exists
	storageDir := "./storage/private/excels"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Generate file path
	filename := req.Filename
	if !strings.HasSuffix(filename, ".xlsx") {
		filename += ".xlsx"
	}

	filePath := filepath.Join(storageDir, filename)

	// Save the Excel file
	if err := f.SaveAs(filePath); err != nil {
		uc.log.Errorf("Failed to save Excel file %s: %v", filename, err)
		return nil, fmt.Errorf("failed to save Excel file: %w", err)
	}

	uc.log.Infof("Successfully exported Excel file: %s with %d rows", filename, len(dataRows))

	// Return response
	response := &entity.ExcelResponse{
		Filename: filename,
		URL:      fmt.Sprintf("/storage/private/excels/%s", filename),
		Rows:     len(dataRows),
	}

	return response, nil
}

// ExportCustomQueryToExcel executes a custom SQL query and exports the results to Excel
func (uc *excelUsecase) ExportCustomQueryToExcel(db *sqlx.DB, req entity.CustomQueryRequest) (*entity.ExcelResponse, error) {
	uc.log.Infof("Starting custom query export with filename: %s", req.Filename)

	// Security check: Only allow SELECT statements
	trimmedQuery := strings.TrimSpace(strings.ToUpper(req.Query))
	if !strings.HasPrefix(trimmedQuery, "SELECT") {
		uc.log.Warnf("Security violation attempt: Non-SELECT query detected: %s", trimmedQuery)
		return nil, fmt.Errorf("only SELECT statements are allowed for security reasons")
	}

	// Prepare query arguments
	var args []interface{}
	if req.Params != nil {
		// Convert map to ordered slice based on parameter names in query
		// This is a simplified approach - in production, you might want to use a more robust solution
		for key, value := range req.Params {
			// Replace placeholders in query (simple approach)
			req.Query = strings.ReplaceAll(req.Query, ":"+key, "?")
			args = append(args, value)
		}
	}

	// Execute the query
	rows, err := db.Queryx(req.Query, args...)
	if err != nil {
		uc.log.Errorf("Failed to execute custom query: %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %w", err)
	}

	// Process rows
	var resultRows []entity.ExcelRow
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a map for the row
		rowMap := make(entity.ExcelRow)
		for i, col := range columns {
			// Handle nil values
			if values[i] == nil {
				rowMap[col] = ""
			} else {
				rowMap[col] = values[i]
			}
		}

		resultRows = append(resultRows, rowMap)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		uc.log.Errorf("Error during row iteration: %v", err)
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	uc.log.Infof("Successfully executed custom query and retrieved %d rows", len(resultRows))

	// Convert query result to Excel data format
	excelData := &entity.ExcelData{
		Headers: columns,
		Rows:    resultRows,
	}

	// Create Excel export request
	exportReq := entity.ExcelExportRequest{
		Data:     excelData,
		Sheet:    req.Sheet,
		Filename: req.Filename,
		Table:    "custom_query", // This is just a placeholder
	}

	// Export to Excel using existing export function
	return uc.ExportExcel(exportReq)
}

// ExportQueryResultToExcel exports a pre-executed query result to Excel
func (uc *excelUsecase) ExportQueryResultToExcel(queryResult *entity.QueryResult, sheet, filename string) (*entity.ExcelResponse, error) {
	uc.log.Infof("Starting export of pre-executed query result to Excel with filename: %s", filename)

	// Convert query result to Excel data format
	// Convert []map[string]interface{} to []ExcelRow
	var excelRows []entity.ExcelRow
	for _, row := range queryResult.Rows {
		excelRow := make(entity.ExcelRow)
		for key, value := range row {
			excelRow[key] = value
		}
		excelRows = append(excelRows, excelRow)
	}

	excelData := &entity.ExcelData{
		Headers: queryResult.Columns,
		Rows:    excelRows,
	}

	// Create Excel export request
	exportReq := entity.ExcelExportRequest{
		Data:     excelData,
		Sheet:    sheet,
		Filename: filename,
		Table:    "custom_query", // This is just a placeholder
	}

	uc.log.Infof("Successfully converted query result to Excel data format with %d rows", len(excelRows))

	// Export to Excel using existing export function
	return uc.ExportExcel(exportReq)
}
