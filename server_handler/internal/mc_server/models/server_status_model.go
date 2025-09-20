package models

// HealthStatus represents the health status of the server
type HealthStatus string

const (
	HealthHealthy  HealthStatus = "HEALTHY"
	HealthSick     HealthStatus = "SICK"
	HealthCritical HealthStatus = "CRITICAL"
	HealthUnknown  HealthStatus = "UNKNOWN"
)

// ServerStatusEnum represents the operational status of the server
type ServerStatusEnum string

const (
	StatusOnline      ServerStatusEnum = "ONLINE"
	StatusOffline     ServerStatusEnum = "OFFLINE"
	StatusMaintenance ServerStatusEnum = "MAINTENANCE"
	StatusUnknown     ServerStatusEnum = "UNKNOWN"
)

// ServerStatus represents the current status of the Minecraft server
type ServerStatus struct {
	Health      HealthStatus     `json:"health" binding:"required"`
	Status      ServerStatusEnum `json:"status" binding:"required"`
	Description string           `json:"description" binding:"required"`
	Uptime      string           `json:"uptime" binding:"required"`
	TPS         float64          `json:"tps" binding:"required"`
}

// NewServerStatus creates a new ServerStatus instance
func NewServerStatus(health HealthStatus, status ServerStatusEnum, description, uptime string, tps float64) *ServerStatus {
	return &ServerStatus{
		Health:      health,
		Status:      status,
		Description: description,
		Uptime:      uptime,
		TPS:         tps,
	}
}
