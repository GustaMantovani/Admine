# SERVER HANDLER #
The server handler starts a minecraft server from a compose file.
It listens to pubsub channels for commands to the server. Responds with zerotier id.

## Configuration ##
The server handler can be configured in a .yaml file or environment variables.
### Yaml ###
The file is ~/.config/admine/server.yaml
```
serverName: "name-of-compose-service"
composeDirectory: "/compose/absolute/path.yaml"
host: "pubsub-host-address"
port: "pubsub-port"
pubsub: "pubsub-type"
senderChannel: "channel-that-responds"
consumerChannels:
- "channel1"
- "channel2"
```

### Environment Variables ###
```
SERVER_NAME "channel"
COMPOSE_DIRECTORY "/path"
CONSUMER_CHANNEL "channel1:channel2"
SENDER_CHANNEL "channel"
PUBSUB "pubsub-type"
HOST "pubsub-host-address"
PORT "pubsub-port"
```
