# SERVER HANDLER

The server handler manages Minecraft servers through Docker Compose.
It listens to Redis pubsub channels for server commands and provides a REST API.

## Configuration

The server handler can be configured via YAML file or environment variables.

### YAML
The configuration file is `server_handler_config.yaml` in the project directory.

```yaml
app:
  self_origin_name: "server_handler"
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
  server_up_timeout: "2m"
  server_off_timeout: "1m"
  server_command_exec_timeout: "30s"
  docker:
    compose_path: "/path/to/docker-compose.yaml"
    container_name: "mine_server"
    service_name: "mine_server"
    rcon_address: "127.0.0.1:25575"
    rcon_password: "password"

web_server:
  host: "0.0.0.0"
  port: 3000
```

## API Endpoints
- `GET /info` - Get server information (version, Java, mods, seed, max players)
- `GET /status` - Get server status (health, uptime, TPS, player count)
- `POST /command` - Execute commands on the server

## PubSub Commands
Listens for messages with tags:
- `server_on` - Start the server
- `server_off` - Stop the server gracefully
- `server_down` - Remove server containers
- `restart` - Restart the server
- `command` - Execute server command
