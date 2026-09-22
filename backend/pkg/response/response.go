package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// JSON writes a standardized JSON response
func JSON(w http.ResponseWriter, status int, success bool, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}

// Error writes a standardized error JSON response
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, false, message, nil)
}

// Message writes a standardized message-only JSON response
func Message(w http.ResponseWriter, status int, success bool, message string) {
	JSON(w, status, success, message, nil)
}
