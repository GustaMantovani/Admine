# Admine Developer Guide

## Overview

Admine is a distributed infrastructure management system designed to automate Minecraft server lifecycle operations in containerized environments. The architecture decouples concerns into specialized components that communicate asynchronously through Redis Pub/Sub. Each module is self-contained with its own configuration, dependencies, and testing infrastructure—refer to the module-specific README for technical implementation details.

## Architecture

![Admine](.readme/Admine.png)

### 1. Server Handler (Go)
The `server_handler` component orchestrates Minecraft server lifecycle management through Docker Compose.

### 2. VPN Handler (Rust)
The `vpn_handler` component manages network access by integrating with ZeroTier.

### 3. Discord Bot (Python)
The `bot` component serves as the primary user interface.

### 4. Message Bus (Redis)
The system uses Redis Pub/Sub as a message broker, enabling asynchronous communication between all components. Pub/Sub channels are:

- `server_channel`: Server lifecycle events (up, down, restarting)
- `command_channel`: Command routing and execution results
- `vpn_channel`: Network state changes and member updates

All messages follow a standardized format for predictable serialization and routing:

```json
{
  "origin": "component_name",
  "tags": ["event_tag_1", "event_tag_2"],
  "message": "a content"
}
```

## Project Utilities

The `utils` directory provides supporting infrastructure and tooling:

### Pub/Sub Utilities (`utils/pubsub/redis/`)
Python scripts for creating and sending messages to Redis Pub/Sub for debugging.

### Release Utilities (`utils/releasing/`)
Automation tooling for building deployable artifacts. The `make-release.nu` script orchestrates packaging of all components into deployment-ready archives. The `templates/Admine-Deploy-Pack/` directory provides a pre-configured deployment structure, reducing configuration overhead during initial setup.

### Mock APIs (`utils/mocks/apis/`)
Docker Compose configurations for mocking all APIs.

## Communication Flow

The request-response pattern in Admine follows a consistent flow:

1. **User executes Discord command** → Bot receives via discord.py
2. **Bot validates permissions** → Routes to CommandHandle
3. **CommandHandle publishes to pub/sub** → Specific channel (server_channel, vpn_channel)
4. **Target service subscribes to channel** → Processes command
5. **Service executes operation** → Publishes result to appropriate channel
6. **EventHandle receives result** → Notifies all Discord channels

The asynchronous architecture enables computationally intensive and long-running server operations to execute without maintaining synchronous connections for extended periods. An example is server initialization, which includes world generation stages that would otherwise require keeping connections open unnecessarily.