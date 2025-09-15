package models

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Message string `json:"message" binding:"required"`
}

// NewErrorResponse creates a new ErrorResponse instance
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
	}
}
