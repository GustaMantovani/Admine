### README for ZeroTier Handler

# VPN Handler

The VPN Handler is a Rust-based component of the Admine infrastructure management solution. It is responsible for managing network configurations and member authorizations within a VPN network. This handler automates the process of adding and removing network members.

## Environment variables:

```bash
export RUST_BACKTRACE=full
export PUBSUB_URL=redis://localhost
export PUBSUB_TYPE=Redis
export VPN_API_URL=https://api.zerotier.com/api/v1
export VPN_API_KEY=token
export VPN_NETWORK_ID=id
export COMMAND_CHANNEL=command_channel
export SERVER_CHANNEL=server_channel
export VPN_CHANNEL=vpn_channel
export DB_PATH=./sled
export STORE_TYPE=Sled
export VPN_RETRY_DELAY_MS=10000
export VPN_RETRY_ATTEMPTS=3
```

## Log Configuration
The application uses log4rs for logging. The log configuration is specified in the `config/log4rs.yaml` file. Below is an example configuration:

```
refresh_rate: 30 seconds

appenders:
  stdout:
    kind: console

  file:
    kind: file
    path: "/tmp/zerotier_handler.log"
    encoder:
      pattern: "{d} - {l} - {m}{n}"

root:
  level: info
  appenders:
    - stdout
    - file
```
