package models

type Response struct {
	Success    bool `json:"success"`
	StatusCode int  `json:"status_code"`
	Data       any  `json:"data,omitempty,omitzero"`
}

type ErrorResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

func NewErrorResponse(status int, err error) *ErrorResponse {
	return &ErrorResponse{
		Success:    false,
		StatusCode: status,
		Error:      err.Error(),
	}
}
