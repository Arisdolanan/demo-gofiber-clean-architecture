package response

type HTTPSuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type JSONSucces struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
