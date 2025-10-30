package response

type HTTPErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  []JSONError `json:"errors"`
}

type JSONError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
