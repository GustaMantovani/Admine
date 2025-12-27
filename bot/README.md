# Admine Discord Bot

Discord bot for managing Minecraft servers through the Admine infrastructure. Provides a user-friendly interface for server control via Discord slash commands.

## Available Commands

Type `/` in Discord to see all available commands.

### Server Management
- `/on` - Start the Minecraft server
- `/off` - Stop the Minecraft server  
- `/restart` - Restart the Minecraft server
- `/status` - Get current server status
- `/info` - Get detailed server information

### Server Commands
- `/command <mine_command>` - Execute a Minecraft command on the server
  - Example: `/command say Hello World!`
  - Example: `/command tp player1 player2`

### VPN Management
- `/auth <vpn_id>` - Authorize a member on the VPN network
- `/vpn_id` - Get the VPN network ID
- `/server_ips` - Get server IP addresses in the VPN

### Administration (Admin Only)
- `/adm <user>` - Grant admin privileges to a Discord user
- `/add_channel` - Add current channel to allowed channels list
- `/remove_channel` - Remove current channel from allowed channels list

## Setup

### Prerequisites
- pyenv
- poetry
- git-hooks
- Make

### Configuration

Create your configuration in `bot_config.json`:

```json
{
    "logging": {
        "level": "DEBUG",
        "file": "./bot.log"
    },
    "providers": {
        "messaging": "DISCORD",
        "pubsub": "REDIS",
        "minecraft": "REST",
        "vpn": "REST"
    },
    "discord": {
        "token": "YOUR_DISCORD_BOT_TOKEN",
        "commandprefix": "!mc",
        "administrators": [
            "admin_user_id_1",
            "admin_user_id_2"
        ],
        "channel_ids": [
            "allowed_channel_id_1"
        ]
    },
    "redis": {
        "connectionstring": "localhost:6379"
    },
    "minecraft": {
        "connectionstring": "http://localhost:3000/api/v1/",
        "token": ""
    },
    "vpn": {
        "connectionstring": "http://localhost:9000",
        "token": ""
    }
}
```

### Run the Bot

```bash
# Using Make
make run

# Using Poetry directly
poetry run python src/main.py
```

## Development

### Git Hooks (Recommended)

The project uses pre-commit hooks for code quality enforcement:

```bash
# Install git hooks (automatically runs on commits)
make git-hooks
# or
poetry run pre-commit install
```

This will automatically run:
- **Ruff linting** with auto-fix
- **Ruff formatting** on all staged files

**Before Submitting**:
   ```bash
   make check     # Ensure all checks pass
   make test      # Run full test suite
   ```

## Architecture

### How It Works

The bot uses a layered architecture with abstractions and providers:

**Core Components:**

1. **Abstractions Layer** - Service interfaces in `external/abstractions/`
   - `MessageService`: Discord communication
   - `PubSubService`: Redis pub/sub for async events
   - `MinecraftServerService`: REST API for server control
   - `VpnService`: VPN network management

2. **Providers Layer** - Concrete implementations using factory pattern
   - Swappable implementations without changing core logic
   - Configure providers via `bot_config.json`

3. **Handlers** - Command and event processing
   - `CommandHandle`: Processes user commands, validates admin permissions, publishes to pub/sub
   - `EventHandle`: Listens to server events, broadcasts notifications to Discord

### Handlers

Handlers are the orchestrators that contain the business logic for each command and route pub/sub events and message service commands to the methods that execute the corresponding action. The bot uses two handlers:

- `CommandHandle` receives commands from Message Service and routes them to appropriate services. It uses a dictionary that maps command names to handler methods. When a command arrives:

    1. The command string is looked up in the `__HANDLES` dictionary
    2. Permission validation occurs via decorator (`@admin_command`)
    3. The corresponding handler method is invoked with arguments

- `EventHandle` subscribes to the Redis pub/sub channel and listens for system events. When an event arrives:

    1. The event's `tags` list is iterated (one event can trigger multiple handlers)
    2. Each tag is looked up in the `__HANDLES` dictionary
    3. Corresponding handler method is invoked with the full event object
    4. Handler usually calls `__notify_all()` which broadcasts to all registered `MessageService` instances
