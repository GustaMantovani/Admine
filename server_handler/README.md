### README for Server Handler

# Server Handler

The server handler is a Go-based Admine component. He is separate in three parts, the Server Handler, Health Checker and Command Handler.

The Server Handler is responsible for set up the Minecraft Server container and comunicate this via a Redis PubSub.

The Health Checker verify the state of the Minecraft Server container.

The Command Handler receives a command via a Redis PubSub and send it to the Minecraft Server.

# Environment variables:

```bash
export REDIS_SERVER_CHANNEL=server_channel
export COMMAND_CHANNEL=command_channel
export NETWORK_ID=xxxxxxxxxxxx
export REDIS_URL=redis://127.0.0.1/
```

# Server Handler config

```yaml
serverName: "name_server_service_in_the_compose"
composeDirectory: "compose_directory_full_name"
```
