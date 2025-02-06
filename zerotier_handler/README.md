### README for ZeroTier Handler

# ZeroTier Handler

The ZeroTier Handler is a Rust-based component of the Admine infrastructure management solution. It is responsible for managing network configurations and member authorizations within a ZeroTier network. This handler automates the process of adding and removing network members.

## Environment variables:

```bash
RUST_BACKTRACE=full
PUBSUB_URL=redis://localhost
PUBSUB_TYPE=Redis
VPN_API_URL=https://api.zerotier.com/api/v1
VPN_API_KEY=token
VPN_NETWORK_ID=id
COMMAND_CHANNEL=command_channel
SERVER_CHANNEL=server_channel
VPN_CHANNEL=vpn_channel
DB_PATH=./sled
STORE_TYPE=Sled
VPN_RETRY_DELAY_MS=10000
VPN_RETRY_ATTEMPTS=3
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