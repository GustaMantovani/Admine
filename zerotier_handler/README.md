### README for ZeroTier Handler

# ZeroTier Handler

The ZeroTier Handler is a Rust-based component of the Admine infrastructure management solution. It is responsible for managing network configurations and member authorizations within a ZeroTier network. This handler automates the process of adding and removing network members.

## Environment variables:

```bash
export RUST_LOG=info
export ZEROTIER_API_BASE_URL=https://api.zerotier.com/api/v1
export ZEROTIER_API_TOKEN=your_token
export ZEROTIER_NETWORK_ID=your_network_id
export ZEROTIER_HANDLER_RETRY_COUNT=3
export ZEROTIER_HANDLER_RETRY_INTERVAL=5
export RECORD_FILE_PATH=record.json
export REDIS_URL=redis://127.0.0.1/
export REDIS_SERVER_CHANNEL=server_channel
export REDIS_COMMAND_CHANNEL=command_channel
export REDIS_VPN_CHANNEL=vpn_channel
```

## Log Configuration
The application uses log4rs for logging. The log configuration is specified in the `log4rs.yaml` file. Below is an example configuration:

```
refresh_rate: 30 seconds

appenders:
  stdout:
    kind: console

  file:
    kind: file
    path: "zerotier_handler.log"
    encoder:
      pattern: "{d} - {l} - {m}{n}"

root:
  level: info
  appenders:
    - stdout
    - file
```