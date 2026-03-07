package models

// ModInstallRequest represents a request to install a mod via URL
type ModInstallRequest struct {
	URL string `json:"url" binding:"required"`
}

// ModInstallResponse represents the accepted response (202)
type ModInstallResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewModInstallResponse creates a new ModInstallResponse
func NewModInstallResponse(status string, message string) *ModInstallResponse {
	return &ModInstallResponse{
		Status:  status,
		Message: message,
	}
}
