package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIErrorResponse struct {
	Code string
	Message string
	RequestID string
}

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *APIErrorResponse `json:"error,omitempty"`
}

// WriteJSON writes a standardized JSON response.
func WriteJSON(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)

	if err := encoder.Encode(response); err != nil {
		// At this point headers are already written.
		// Only log the failure.
		log.Printf("failed to encode JSON response: %v", err)
	}
}

// SuccessResponse writes a successful JSON response.
func SuccessResponse(w http.ResponseWriter, statusCode int, data any) {
	WriteJSON(w, statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse writes an error JSON response.
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, APIResponse{
		Success: false,
		Error:   &APIErrorResponse {
			Message: message,
		},
	})
}