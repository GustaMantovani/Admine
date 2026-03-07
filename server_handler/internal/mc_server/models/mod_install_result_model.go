package models

// ModInstallResult represents the result of a mod installation operation
type ModInstallResult struct {
	FileName string `json:"file_name"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// NewModInstallResult creates a new ModInstallResult
func NewModInstallResult(fileName string, success bool, message string) *ModInstallResult {
	return &ModInstallResult{
		FileName: fileName,
		Success:  success,
		Message:  message,
	}
}
