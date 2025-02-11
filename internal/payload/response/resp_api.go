package response

type SuccessResponse struct {
	ResponseStatus  bool        `json:"status"`
	ResponseMessage string      `json:"message"`
	Data            interface{} `json:"data,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message any    `json:"message"`
}

type ErrorResponse struct {
	ResponseStatus  bool          `json:"status"`
	ResponseMessage string        `json:"message"`
	Errors          []ErrorDetail `json:"errors,omitempty"`
}
