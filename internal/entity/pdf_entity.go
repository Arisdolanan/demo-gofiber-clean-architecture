package entity

// PDFRequest represents the request structure for PDF generation
type PDFRequest struct {
	HTMLContent string `json:"html_content" validate:"required"`
	Filename    string `json:"filename" validate:"required"`
	Title       string `json:"title,omitempty"`
	Author      string `json:"author,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
}

// PDFResponse represents the response structure for PDF generation
type PDFResponse struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
}
