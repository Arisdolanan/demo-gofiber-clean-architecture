package entity

// ExcelImportRequest represents the request structure for importing Excel data
type ExcelImportRequest struct {
	FileData []byte `json:"file_data" validate:"required"`
	Sheet    string `json:"sheet" validate:"required"`
	Table    string `json:"table" validate:"required"`
}

// ExcelExportRequest represents the request structure for exporting data to Excel
type ExcelExportRequest struct {
	Data     interface{} `json:"data" validate:"required"`
	Sheet    string      `json:"sheet" validate:"required"`
	Filename string      `json:"filename" validate:"required"`
	Table    string      `json:"table" validate:"required"`
}

// CustomQueryRequest represents the request structure for executing custom queries
type CustomQueryRequest struct {
	Query    string                 `json:"query" validate:"required"`
	Params   map[string]interface{} `json:"params,omitempty"`
	Sheet    string                 `json:"sheet" validate:"required"`
	Filename string                 `json:"filename" validate:"required"`
}

// QueryResult represents the result of a custom query execution
type QueryResult struct {
	Columns []string                 `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
}

// ExcelResponse represents the response structure for Excel operations
type ExcelResponse struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Rows     int    `json:"rows"`
}

// ExcelRow represents a row of Excel data
type ExcelRow map[string]interface{}

// ExcelData represents Excel data structure
type ExcelData struct {
	Headers []string   `json:"headers"`
	Rows    []ExcelRow `json:"rows"`
}
