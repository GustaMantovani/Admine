# Admine - Infrastructure Manager for Minecraft Servers

Admine is a comprehensive infrastructure management solution for Minecraft servers running on Linux systems. It provides automated server lifecycle management, VPN connectivity through ZeroTier, and Discord-based administration interface.

## What is Admine?

Admine automates the complete management of Minecraft servers in a containerized environment, solving the common problem of running servers behind NAT (Network Address Translation) or in private networks where direct public access isn't available. The system is designed for server administrators who want to host Minecraft servers in private networks (home networks, cloud instances behind NAT, etc.) while providing secure and easy access to players through VPN connectivity.

## System Architecture

![Admine](.readme/Admine.png)

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
