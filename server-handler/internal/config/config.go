package config

import (
	"sync"
)

type Config struct {
	ComposeAbsPath string
}

var instance *Config
var once sync.Once

// Obter instancia Singleton da configuração do servidor
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{ComposeAbsPath: "/home/andre/pgm/pessoal/Admine/minecraft-server/docker-compose.yaml"}
	})

	return instance
}
