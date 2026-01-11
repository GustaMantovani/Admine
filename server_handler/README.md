# SERVER HANDLER

The server handler manages Minecraft servers through Docker Compose.
It listens to Redis pubsub channels for server commands and provides a REST API.

## Configuration

### Default Values
The server handler comes with sensible defaults for all configuration fields. If the configuration file doesn't exist or if specific fields are missing, the following defaults will be used:

- **App**
  - `self_origin_name`: "server"
  - `log_file_path`: "/tmp/server_handler.log"
  - `log_level`: "INFO"

- **PubSub**
  - `type`: "redis"
  - Redis: `addr`: "localhost:6379", `password`: "", `db`: 0
  - Channels: server_channel, command_channel, vpn_channel

- **Minecraft Server**
  - `runtime_type`: "docker"
  - `server_type`: "fabric"
  - `server_up_timeout`: 2m
  - `server_off_timeout`: 1m
  - `server_command_exec_timeout`: 30s
  - `rcon_address`: "127.0.0.1:25575"
  - `rcon_password`: "admineRconPassword!"
  - Docker: `compose_path`: auto-generated based on server_type, `container_name`: "minecraft_server", `service_name`: "minecraft_server"

- **Web Server**
  - `host`: "0.0.0.0"
  - `port`: 3000

### YAML
The configuration file is `server_handler_config.yaml` in the project directory. All fields are optional and will use defaults if not specified.

```yaml
app:
  self_origin_name: "server"
  log_file_path: "/tmp/server_handler.log"
  log_level: "DEBUG"

pubsub:
  type: "redis"
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
  admine_channels_map:
    server_channel: "server_channel"
    command_channel: "command_channel"
    vpn_channel: "vpn_channel"

minecraft_server:
  runtime_type: "docker"
  server_type: "fabric"  # Used to auto-generate compose_path if not specified
  server_up_timeout: "2m"
  server_off_timeout: "1m"
  server_command_exec_timeout: "30s"
  rcon_address: "127.0.0.1:25575"
  rcon_password: "123456"
  docker:
    compose_path: "../minecraft_server/fabric/docker-compose.yaml"  # Optional: auto-generated from server_type if omitted
    container_name: "mine_server"
    service_name: "mine_server"

web_server:
  host: "0.0.0.0"
  port: 3000
```

**Note:** If `compose_path` is not specified, it will be automatically generated as `../minecraft_server/{server_type}/docker-compose.yaml`.
