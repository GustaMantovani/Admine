package mcserver

type MinecraftServer interface {
	Start() error
	Stop() error
	Down() error
	Restart() error
	Status() (string, error)
	Info() (string, error)
	StartUpInfo() string
	ExecuteCommand(command string) (string, error)
}
