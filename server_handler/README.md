# SERVER HANDLER #
The server handler up minecraft server from a compose file.
It listens to pubsub channels for commands to the server. Responds with zerotier id.

## Configuration ##
The server handler can be configured in a .yaml file or env vars.
### Yaml ###
The file is ~/.config/server.yaml
```
serverName: "name-of-compose-service"
composeDirectory: "/compose/absolute/path.yaml"
host: "pubsub-host-adress"
port: "pubsub-port"
pubsub: "pubsub-type"
senderChannel: "channel-that-responds"
consummerChannels:
- "channel1"
- "channel2"
```

### Env ###
```
SERVER_NAME "channel"
COMPOSE_DIRECTORY "/path"
CONSUMER_CHANNEL "channel1:channel2"
SENDER_CHANNEL "channel"
PUBSUB "pubsub-type"
HOST "pubsub-host-adress"
PORT "pubsub-port"
```
