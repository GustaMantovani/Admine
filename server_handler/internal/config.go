package internal

type Config struct {
	// App identity
	SelfOriginName string `yaml:"self_origin_name"`

	// Redis configuration
	RedisAddr     string `yaml:"redis_addr"`
	RedisPassword string `yaml:"redis_password"`
	RedisDB       int    `yaml:"redis_db"`

	// Docker Compose
	DockerComposePath string `yaml:"compose_path"`

	// Other optional fields
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// LoadConfig reads a YAML file and unmarshals into Config
func LoadConfig(path string) (*Config, error) {
	// exemplo usando gopkg.in/yaml.v3
	// f, err := os.Open(path)
	// ...
	// yaml.NewDecoder(f).Decode(&cfg)
	return &Config{}, nil
}
