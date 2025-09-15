package models

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
	return &CommandResult{
		Output: output,
	}
}

// NewCommandResultWithExitCode creates a new CommandResult with output and exit code
func NewCommandResultWithExitCode(output string, exitCode int) *CommandResult {
	return &CommandResult{
		ExitCode: &exitCode,
		Output:   output,
	}
}
