package models

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
