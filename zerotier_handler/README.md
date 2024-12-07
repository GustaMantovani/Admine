### README for ZeroTier Handler

# ZeroTier Handler

The ZeroTier Handler is a Rust-based component of the Admine infrastructure management solution. It is responsible for managing network configurations and member authorizations within a ZeroTier network. This handler automates the process of adding and removing network members, ensuring that the Minecraft server is always accessible via a consistent IP address.

## Features

- **ZeroTier API Integration**: Manages network members using the ZeroTier API.
- **IP Management**: Automates IP address assignments for the Minecraft server.
- **State Persistence**: Stores network member states for recovery and consistency.
- **Redis Integration**: Communicates with other components via Redis Pub/Sub channels.

## Environment variables:

```bash
export ZEROTIER_API_BASE_URL=https://my.zerotier.com/api
export ZEROTIER_API_TOKEN=your_token
export ZEROTIER_NETWORK_ID=your_network_id
export RECORD_FILE_PATH=/path/to/record/file
export REDIS_URL=redis://localhost:6379
export REDIS_SERVER_CHANNEL=server_channel
export REDIS_COMMAND_CHANNEL=command_channel
export REDIS_VPN_CHANNEL=vpn_channel
```
