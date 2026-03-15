package config

import (
	"errors"
	"fmt"
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
	RuntimeType              string                `yaml:"runtime_type"`
	Docker                   DockerConfig          `yaml:"docker"`
	Image                    MinecraftImageConfig  `yaml:"image"`
	ZeroTier                 ZeroTierSidecarConfig `yaml:"zerotier"`
	ServerOnTimeout          time.Duration         `yaml:"server_up_timeout"`
	ServerOffTimeout         time.Duration         `yaml:"server_off_timeout"`
	ServerCommandExecTimeout time.Duration         `yaml:"server_command_exec_timeout"`
	ModInstallTimeout        time.Duration         `yaml:"mod_install_timeout"`
	RconAddress              string                `yaml:"rcon_address"`
	RconPassword             string                `yaml:"rcon_password"`
}

// DockerConfig holds Docker-specific runtime settings.
type DockerConfig struct {
	// ComposeOutputPath is where the generated docker-compose.yaml is written.
	ComposeOutputPath string `yaml:"compose_output_path"`
	ContainerName     string `yaml:"container_name"`
	ServiceName       string `yaml:"service_name"`
	// DataPath is the host directory mounted as /data inside the itzg container.
	DataPath string `yaml:"data_path"`
}

// MinecraftImageConfig controls the itzg/docker-minecraft-server image behaviour.
type MinecraftImageConfig struct {
	// Type is the server type passed as the TYPE env var (e.g. FABRIC, FORGE, PAPER, VANILLA).
	Type string `yaml:"type"`
	// Version is the Minecraft version (e.g. "1.20.1").
	Version string `yaml:"version"`
	// FabricLoaderVersion is passed as FABRIC_LOADER_VERSION when non-empty.
	FabricLoaderVersion string `yaml:"fabric_loader_version"`
	// ForgeVersion is passed as FORGE_VERSION when non-empty.
	ForgeVersion string `yaml:"forge_version"`
	// ModpackURL is passed as MODPACK when non-empty (URL to a modpack archive).
	ModpackURL string `yaml:"modpack_url"`
	// JavaVersion sets the image tag used to select the JDK (e.g. "java21", "java17").
	JavaVersion string `yaml:"java_version"`
	// ExtraEnv is a map of additional environment variables forwarded verbatim to the itzg image.
	ExtraEnv map[string]string `yaml:"extra_env"`
}

// ZeroTierSidecarConfig controls the optional ZeroTier sidecar container.
type ZeroTierSidecarConfig struct {
	// Enabled controls whether the ZeroTier sidecar service is included in the generated compose file.
	Enabled bool `yaml:"enabled"`
	// NetworkID is the ZeroTier network to join.
	NetworkID string `yaml:"network_id"`
	// ContainerName is the name assigned to the ZeroTier container.
	ContainerName string `yaml:"container_name"`
	// ApiSecret is optionally passed as ZEROTIER_API_SECRET inside the ZeroTier container.
	ApiSecret string `yaml:"api_secret"`
}

type WebServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// NewDefaultConfig returns a Config with sensible defaults
func NewDefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			SelfOriginName: "server",
			LogFilePath:    "/tmp/admine/logs/server_handler.log",
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
			ModInstallTimeout:        2 * time.Minute,
			RconAddress:              "127.0.0.1:25575",
			RconPassword:             "admineRconPassword!",
			Docker: DockerConfig{
				ComposeOutputPath: "./generated/docker-compose.yaml",
				ContainerName:     "mine_server",
				ServiceName:       "mine_server",
				DataPath:          "./minecraft-data",
			},
			Image: MinecraftImageConfig{
				Type:    "FABRIC",
				Version: "1.20.1",
			},
			ZeroTier: ZeroTierSidecarConfig{
				Enabled:       true,
				ContainerName: "zerotier",
			},
		},
		WebSever: WebServerConfig{
			Host: "0.0.0.0",
			Port: 3000,
		},
	}
}

// LoadConfig reads a YAML file into Config, using defaults for missing values
func LoadConfig(path string) (*Config, error) {
	cfg := NewDefaultConfig()

	fmt.Printf("Default config created: %+v\n", cfg)

	exists, err := pathExists(path)
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

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return cfg, nil
	}

	if cfg.MinecraftServer.Docker.ComposeOutputPath == "" {
		cfg.MinecraftServer.Docker.ComposeOutputPath = "./generated/docker-compose.yaml"
	}

	fmt.Printf("New config loaded: %+v\n", cfg)
	return cfg, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
