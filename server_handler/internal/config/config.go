package config

import (
	"fmt"
	"os"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/pkg"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App AppConfig `yaml:"app"`

	PubSub PubSubConfig `yaml:"pubsub"`

	MinecraftServer MinecraftServerConfig `yaml:"minecraft_server"`

	WebSever WebServerConfig `yaml:"web_server"`
}

type AppConfig struct {
	SelfOriginName string `yaml:"self_origin_name"`
	LogFilePath    string `yaml:"log_file_path"`
	LogLevel       string `yaml:"log_level"`
}

type PubSubConfig struct {
	Type              string            `yaml:"type"`
	Redis             RedisConfig       `yaml:"redis"`
	AdmineChannelsMap AdmineChannelsMap `yaml:"admine_channels_map"`
}

type AdmineChannelsMap struct {
	ServerChannel  string `yaml:"server_channel"`
	CommandChannel string `yaml:"command_channel"`
	VpnChannel     string `yaml:"vpn_channel"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type MinecraftServerConfig struct {
	RuntimeType              string        `yaml:"runtime_type"`
	Docker                   DockerConfig  `yaml:"docker"`
	ServerOnTimeout          time.Duration `yaml:"server_up_timeout"`
	ServerOffTimeout         time.Duration `yaml:"server_off_timeout"`
	ServerCommandExecTimeout time.Duration `yaml:"server_command_exec_timeout"`
	RconAddress              string        `yaml:"rcon_address"`
	RconPassword             string        `yaml:"rcon_password"`
}

type DockerConfig struct {
	ComposePath   string `yaml:"compose_path"`
	ContainerName string `yaml:"container_name"`
	ServiceName   string `yaml:"service_name"`
}

type WebServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// NewDefaultConfig returns a Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			SelfOriginName: "server",
			LogFilePath:    "/tmp/server_handler.log",
			LogLevel:       "INFO",
		},
		PubSub: PubSubConfig{
			Type: "redis",
			Redis: RedisConfig{
				Addr:     "localhost:6379",
				Password: "",
				Db:       0,
			},
			AdmineChannelsMap: AdmineChannelsMap{
				ServerChannel:  "server_channel",
				CommandChannel: "command_channel",
				VpnChannel:     "vpn_channel",
			},
		},
		MinecraftServer: MinecraftServerConfig{
			RuntimeType:              "docker",
			ServerOnTimeout:          2 * time.Minute,
			ServerOffTimeout:         1 * time.Minute,
			ServerCommandExecTimeout: 30 * time.Second,
			RconAddress:              "127.0.0.1:25575",
			RconPassword:             "admineRconPassword!",
			Docker: DockerConfig{
				ComposePath:   "./docker-compose.yaml",
				ContainerName: "minecraft_server",
				ServiceName:   "minecraft_server",
			},
		},
		WebSever: WebServerConfig{
			Host: "0.0.0.0",
			Port: 3000,
		},
	}
}

// LoadConfig reads YAML file into Config with default values
func LoadConfig(path string) (*Config, error) {
	// Start with default configuration
	cfg := NewDefaultConfig()

	fmt.Printf("Default config created: %+v\n", cfg)

	exists, err := pkg.PathExists(path)

	if err != nil {
		return nil, err
	}

	if !exists {
		return cfg, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Override defaults with values from YAML file
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return cfg, nil
	}

	fmt.Printf("New config loaded: %+v\n", cfg)
	return cfg, nil
}
