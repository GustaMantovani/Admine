package models

// LogsResponse represents a logs response from the API
type LogsResponse struct {
	Lines []string `json:"lines"`
	Total int      `json:"total"`
}

// NewLogsResponse creates a new LogsResponse instance
func NewLogsResponse(lines []string) *LogsResponse {
	return &LogsResponse{
		Lines: lines,
		Total: len(lines),
	}
}
