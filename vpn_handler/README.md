### README for ZeroTier Handler

# VPN Handler

Default values are provided, so the application can start without a full configuration. Still, there is a minimal configuration that must be provided for the VPN integration to work correctly: `vpn_config.api_key` and `vpn_config.network_id`. You can override any other entries individually or provide the full configuration as you prefer. The VPN Handler is a Rust-based component of the Admine. It is responsible for managing network configurations and member authorizations within a VPN network. This handler automates the process of adding and removing network members.

## Configuration file:

The configuration is specified in the `./etc/vpn_handler_config.toml` file. Below is an example configuration (note the required minimal fields under `vpn_config`):

```toml
self_origin_name = "vpn"

[api_config]
host = "localhost"
port = 9000

[pub_sub_config]
url = "redis://localhost:6379"
pub_sub_type = "Redis"

[vpn_config]
api_url = "https://api.zerotier.com/api/v1"
api_key = "your_api_key_here"
network_id = "your_network_id_here"
vpn_type = "Zerotier"

[db_config]
path = "./etc/sled/vpn_store.db"
store_type = "Sled"

[admine_channels_map]
server_channel = "server_channel"
command_channel = "command_channel"
vpn_channel = "vpn_channel"

[retry_config]
attempts = 5
delay = { secs = 3, nanos = 0 }
```

## Log Configuration

The application uses log4rs for logging. No configuration is required for logs to work and for the application to start, since default values already exist. Still, you can define individual configuration entries or the entire configuration as you prefer. The log configuration is specified in the `./etc/log4rs.yaml` file. Below is an example configuration:

```
refresh_rate: 30 seconds

appenders:
  stdout:
    kind: console

  file:
    kind: file
    path: "./etc/zerotier_handler.log"
    encoder:
      pattern: "{d} - {l} - {m}{n}"

root:
  level: info
  appenders:
    - stdout
    - file
```
