package models

import (
	"encoding/json"
	"net/http"
)

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

// any e interface{} son lo mismo , any es solo syntaxis sugar para interface{}, any se introdujo en GO 1.18
func ResponseWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
