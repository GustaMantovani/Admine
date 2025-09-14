package mcserver

type MinecraftServer interface {
	Start() error
	Stop() error
	Restart() error
	Status() (string, error)
	Info() (string, error)
	ExecuteCommand(command string) (string, error)
}
