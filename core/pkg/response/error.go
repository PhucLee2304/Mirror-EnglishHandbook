package response

type ErrorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"`
}

type HTTPError struct {
	StatusCode int
	Error      ErrorResponse
}
