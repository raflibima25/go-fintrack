package model

type ApiResponse struct {
	ResponseStatus  bool        `json:"status"`
	ResponseMessage string      `json:"message"`
	Data            interface{} `json:"data,omitempty"`
}
