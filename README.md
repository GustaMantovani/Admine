## Core Components

### 1. Server Handler (Go)
The `server_handler` component manages Minecraft server containers:

- Starts, stops, and monitors server containers
- Provides REST API endpoints for server control
- Executes Minecraft commands via RCON
- Collects server status information (health, TPS, uptime)
- Detects server type (Vanilla, Forge, Fabric, etc.)

### 2. VPN Handler (Rust)
The `vpn_handler` component manages network connectivity:

- Manages ZeroTier network members and authorizations
- Automates IP address assignments
- Provides an API for VPN operations
- Persists network state in a key-value store

### 3. Discord Bot (Python)
The `bot` component provides the user interface:

- Processes commands from Discord users
- Displays server status and information
- Controls access to administrative commands
- Sends notifications about server events
- Routes commands to appropriate components

### 4. Message Bus (Redis)
The system uses Redis Pub/Sub for communication between components:

- `server_channel`: Server lifecycle events
- `command_channel`: Command routing
- `vpn_channel`: Network configuration updates
## Minecraft Server Support

The `minecraft_server` directory contains Docker configurations for:

- Forge servers
- Fabric servers
- Support for various Minecraft versions

## Getting Started

1. Configure each component using their respective configuration files
2. Start Redis (see utils/pubsub/redis/)
3. Launch the server handler
4. Launch the VPN handler
5. Start the Discord bot
## Command Interface

Users can interact with the system through Discord commands:

- `/on` - Start the server
- `/off` - Stop the server
- `/restart` - Restart the server
- `/status` - Get server status
- `/info` - Get server information
- `/command <cmd>` - Execute Minecraft commands
- `/auth <id>` - Authorize VPN member
- `/vpn_id` - Get VPN network ID
- `/server_ips` - Get server IP addresses
## API Endpoints

### Server Handler
The server handler provides REST API endpoints:

- `GET /api/v1/info` - Get server information
- `GET /api/v1/status` - Get server status
- `POST /api/v1/command` - Execute commands

### VPN Handler
The VPN handler provides REST API endpoints:

- `GET /server-ips` - Get server IP addresses in VPN
- `GET /vpn-id` - Get internal VPN network ID
- `POST /auth-member` - Authorize a member on VPN network

## Development

Each component has its own build and development workflow. See the `README.md` files in each directory for specific instructions.