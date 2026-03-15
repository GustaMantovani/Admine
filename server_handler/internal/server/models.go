package server

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

// ServerInfo represents the Minecraft server information
type ServerInfo struct {
	MinecraftVersion string `json:"minecraftVersion" binding:"required"`
	JavaVersion      string `json:"javaVersion" binding:"required"`
	ModEngine        string `json:"modEngine" binding:"required"`
	MaxPlayers       int    `json:"maxPlayers" binding:"required"`
	Seed             string `json:"seed" binding:"required"`
}

// NewServerInfo creates a new ServerInfo instance
func NewServerInfo(minecraftVersion, javaVersion, modEngine, seed string, maxPlayers int) *ServerInfo {
	return &ServerInfo{
		MinecraftVersion: minecraftVersion,
		JavaVersion:      javaVersion,
		ModEngine:        modEngine,
		MaxPlayers:       maxPlayers,
		Seed:             seed,
	}
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	ExitCode *int   `json:"exitCode,omitempty"`
	Output   string `json:"output" binding:"required"`
}

// NewCommandResult creates a new CommandResult instance
func NewCommandResult(output string, exitCode *int) *CommandResult {
	return &CommandResult{
		ExitCode: exitCode,
		Output:   output,
	}
}

// NewCommandResultWithOutput creates a new CommandResult with only output
func NewCommandResultWithOutput(output string) *CommandResult {
	return &CommandResult{Output: output}
}

// NewCommandResultWithExitCode creates a new CommandResult with output and exit code
func NewCommandResultWithExitCode(output string, exitCode int) *CommandResult {
	return &CommandResult{
		ExitCode: &exitCode,
		Output:   output,
	}
}

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

// ModListResult represents the list of installed mods
type ModListResult struct {
	Mods  []string `json:"mods"`
	Total int      `json:"total"`
}

// NewModListResult creates a new ModListResult
func NewModListResult(mods []string) *ModListResult {
	return &ModListResult{
		Mods:  mods,
		Total: len(mods),
	}
}
