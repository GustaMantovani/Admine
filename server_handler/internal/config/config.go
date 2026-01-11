package config

import (
	"os"
	"time"

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
	RconAddress   string `yaml:"rcon_address"`
	RconPassword  string `yaml:"rcon_password"`
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

// LoadConfig reads YAML file into Config
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
