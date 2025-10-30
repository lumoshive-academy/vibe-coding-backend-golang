package response

import (
	"encoding/json"
	"net/http"
)

// Message represents the base JSON payload returned to clients.
type Message struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// Write marshals the message and writes it to the response writer.
func Write(w http.ResponseWriter, statusCode int, message Message) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(message)
}

// Success constructs a success message payload.
func Success(data interface{}) Message {
	return Message{
		Status: "success",
		Data:   data,
	}
}

// Failure constructs an error message payload.
func Failure(err interface{}) Message {
	return Message{
		Status: "error",
		Error:  err,
	}
}
