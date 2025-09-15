package models

// Command represents a command to be executed on the Minecraft server
type Command struct {
	Command string `json:"command" binding:"required"`
}

// NewCommand creates a new Command instance
func NewCommand(command string) *Command {
	return &Command{
		Command: command,
	}
}
