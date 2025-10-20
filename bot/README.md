# Admine Discord Bot

Discord bot for managing Minecraft servers through the Admine infrastructure. Provides a user-friendly interface for server control, VPN management, and system monitoring.

## Features

- **Server Control**: Start, stop, restart Minecraft servers
- **Command Execution**: Run Minecraft commands remotely
- **VPN Management**: Authorize members and get network information
- **Server Monitoring**: Check server status and information
- **Access Control**: Admin-only commands with role-based permissions
- **Real-time Communication**: Redis Pub/Sub integration

## Prerequisites

- Python 3.12+
- Poetry (dependency management)
- Redis server running
- Discord bot token
- Access to Admine server handler and VPN handler APIs

## Quick Start

### 1. Install Dependencies

```bash
# Install Poetry if not already installed
curl -sSL https://install.python-poetry.org | python3 -

# Complete development setup (installs dependencies + git hooks)
make dev-setup

# Or install dependencies only
make install
# or
poetry install
```

### 2. Git Hooks Setup (Recommended)

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

### 3. Configuration

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

### 4. Run the Bot

```bash
# Using Make
make run

# Using Poetry directly
poetry run python src/main.py

# Development mode with auto-reload
make dev
```

## Available Commands

The bot uses Discord slash commands (type `/` to see available commands).

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

## Configuration Details

### Discord Settings
- **token**: Your Discord bot token from Discord Developer Portal
- **commandprefix**: Legacy setting (not used with slash commands)
- **administrators**: List of Discord user IDs with admin privileges
- **channel_ids**: List of channel IDs where the bot will respond

### Service Connections
- **redis**: Redis server for Pub/Sub communication
- **minecraft**: Server handler API endpoint
- **vpn**: VPN handler API endpoint

## Development

### Project Structure
```
src/
├── main.py              # Entry point
└── bot/
    ├── bot.py           # Main bot class
    ├── config.py        # Configuration management
    ├── logger.py        # Logging utilities
    ├── exceptions.py    # Custom exceptions
    ├── handles/         # Command and event handlers
    ├── external/        # External service integrations
    └── models/          # Data models
```

### Available Make Commands

```bash
make help           # Show all available commands
make install        # Install dependencies
make run           # Run the bot
make dev           # Run in development mode
make test          # Run tests
make lint          # Run linting (ruff check)
make lint-fix      # Run linting with auto-fix
make format        # Format code with ruff
make format-check  # Check if code is formatted correctly
make check         # Run all checks (lint + format-check)
make fix           # Fix all auto-fixable issues (lint-fix + format)
make clean         # Clean cache files
make logs          # Show bot logs
make git-hooks     # Install pre-commit git hooks
make dev-setup     # Complete development setup
make env-info      # Show environment information
```

### Code Quality

The project uses strict code quality enforcement:

- **Ruff**: For linting and formatting (replaces Black + Flake8)
- **MyPy**: For type checking
- **Pytest**: For testing
- **Pre-commit hooks**: Automatically enforces quality on commits

#### Ruff Configuration
- **Line length**: 120 characters
- **Target**: Python 3.12
- **Rules**: Pyflakes (F), pycodestyle (E, W), isort (I)
- **Auto-fix**: Enabled for most rules
- **Formatting**: Double quotes, 4-space indentation

#### Development Workflow
```bash
# Before committing (or use git hooks)
make fix           # Fix all auto-fixable issues
make check         # Verify all checks pass

# Individual operations
make lint          # Check for linting issues
make format        # Format code
make test          # Run tests

# Check environment
make env-info      # Show Python/Poetry versions
```

## Logging

Logs are written to `/tmp/bot.log` by default. The logging configuration includes:
- Info level and above to console
- All levels to file
- Structured logging with timestamps

## Error Handling

The bot includes comprehensive error handling:
- Discord API errors
- Network connectivity issues
- Redis connection problems
- Invalid command syntax
- Permission violations

## Security Notes

1. **Token Security**: Keep your Discord bot token secure and never commit it to version control
2. **Admin Access**: Carefully manage administrator user IDs
3. **Channel Restrictions**: Use channel_ids to limit bot access to specific channels
4. **Network Access**: Ensure Redis and API endpoints are properly secured

## Troubleshooting

### Common Issues

1. **Bot not responding**:
   - Check Discord token is valid
   - Verify bot has necessary permissions in Discord server
   - Ensure bot is added to the correct channels

2. **Commands failing**:
   - Check Redis connection
   - Verify server handler and VPN handler are running
   - Check API endpoints in configuration

3. **Permission errors**:
   - Verify user ID is in administrators list
   - Check channel ID is in allowed channels

### Debug Mode

Enable debug logging by modifying the logger configuration in `src/main.py`.

## Contributing

1. **Setup Development Environment**:
   ```bash
   make dev-setup  # Installs dependencies + git hooks
   ```

2. **Follow Code Quality Standards**:
   - Git hooks will automatically run on commits
   - Manual checks: `make check`
   - Auto-fix issues: `make fix`

3. **Development Workflow**:
   - Write code following existing patterns
   - Add tests for new features
   - Update documentation as needed
   - Commit with descriptive messages

4. **Before Submitting**:
   ```bash
   make check     # Ensure all checks pass
   make test      # Run full test suite
   ```

The pre-commit hooks will automatically handle code formatting and basic linting, making the development process smoother.
