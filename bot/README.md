# Discord Bot

The bot is the user-facing interface for Admine. It:

- Receives Discord slash commands and translates them into Redis Pub/Sub messages or direct REST API calls
- Subscribes to Pub/Sub channels and relays server events back to Discord channels
- Manages bot configuration (admins, authorized channels) at runtime via Discord commands

---

## Package structure

```
bot/
├── src/
│   ├── main.py                          # Entrypoint: config, logging, signal handling
│   └── bot/
│       ├── bot.py                       # Bot lifecycle: constructs all services, starts tasks
│       ├── config.py                    # JSON config loader with deep-merge defaults
│       ├── logger.py                    # Loguru setup
│       ├── exceptions.py                # ConfigError, ConfigFileError
│       ├── models/                      # Pydantic data models
│       │   ├── admine_message.py        # AdmineMessage envelope
│       │   ├── minecraft_server_info.py # /info response model
│       │   ├── minecraft_server_status.py # /status response model
│       │   ├── resource_usage.py        # /resources response model
│       │   └── logs_response.py         # /logs response model
│       ├── handles/
│       │   ├── command_handle.py        # Routes Discord commands → services
│       │   └── event_handle.py          # Routes Pub/Sub events → Discord notifications
│       └── services/
│           ├── messaging/
│           │   ├── message_service.py           # MessageService ABC
│           │   └── discord_message_service.py   # discord.py implementation + factory
│           ├── minecraft/
│           │   ├── minecraft_server_service.py  # MinecraftServerService ABC
│           │   └── server_handler_api_service.py # server_handler REST API client + factory
│           ├── pubsub/
│           │   ├── pubsub_service.py             # PubSubService ABC
│           │   └── redis_pubsub_service.py       # Redis implementation + factory
│           └── vpn/
│               ├── vpn_service.py                # VpnService ABC
│               └── api_vpn_service.py            # vpn_handler REST API client + factory
└── tests/
    └── ...                              # pytest test suite (unittest.mock)
```

---

## Dependency flow

```
main.py
  └─ Config(bot_config.json)
  └─ Bot(config)
       └─ PubSubServiceFactory.create(...)      ─► RedisPubSubService
       └─ MinecraftServiceFactory.create(...)   ─► ServerHandlerApiService
       └─ VpnServiceFactory.create(...)         ─► ApiVpnService
       └─ CommandHandle(pubsub, minecraft, vpn, config)
       └─ EventHandle(message_services)
       └─ MessageServiceFactory.create(...)     ─► DiscordMessageService
  └─ bot.start()
       └─ message_service.connect()             ─► discord.py event loop (task)
       └─ pubsub_service.listen_message(...)    ─► Redis subscription loop (task)
```

The two async tasks run concurrently via `asyncio.gather`. Shutdown cancels both tasks, closes the Redis connection, and disconnects the Discord client.

---

## Service abstractions

Each service domain has an ABC (Abstract Base Class) that the concrete implementation fulfills. Factories select the implementation based on the config's `providers` section, making it possible to swap implementations without touching handlers.

| ABC | Config key | Implementation |
|---|---|---|
| `MessageService` | `providers.messaging` | `DiscordMessageService` (`DISCORD`) |
| `PubSubService` | `providers.pubsub` | `RedisPubSubService` (`REDIS`) |
| `MinecraftServerService` | `providers.minecraft` | `ServerHandlerApiService` (`REST`) |
| `VpnService` | `providers.vpn` | `ApiVpnService` (`REST`) |

---

## CommandHandle

`bot/handles/command_handle.py` — receives a command name + args from the message service and dispatches to a handler method via a string-keyed dictionary.

Permission model: handlers decorated with `@admin_command` require the calling user's ID to be in `discord.administrators`. If not, the handler returns `"Unauthorized command usage"` without executing.

| Command | Admin only | Action |
|---|---|---|
| `on` | yes | Publishes `AdmineMessage(tags=["server_on"])` to `server_channel` |
| `off` | yes | Publishes `AdmineMessage(tags=["server_off"])` to `server_channel` |
| `restart` | yes | Publishes `AdmineMessage(tags=["restart"])` to `server_channel` |
| `command` | yes | Calls `MinecraftServerService.command()` (REST POST /command) |
| `info` | no | Calls `MinecraftServerService.get_info()` (REST GET /info) |
| `status` | no | Calls `MinecraftServerService.get_status()` (REST GET /status) |
| `resources` | yes | Calls `MinecraftServerService.get_resources()` (REST GET /resources) |
| `logs` | yes | Calls `MinecraftServerService.get_logs(n)` (REST GET /logs?n=N) |
| `install_mod` | yes | Calls `MinecraftServerService.install_mod_url()` or `install_mod_file()` |
| `list_mods` | yes | Calls `MinecraftServerService.list_mods()` |
| `remove_mod` | yes | Calls `MinecraftServerService.remove_mod(filename)` |
| `auth` | no | Calls `VpnService.auth_member(id)` (REST POST /auth-member) |
| `vpn_id` | no | Calls `VpnService.get_vpn_id()` (REST GET /vpn-id) |
| `server_ips` | no | Calls `VpnService.get_server_ips()` (REST GET /server-ips) |
| `adm` | yes | Appends user ID to `discord.administrators` in config and saves |
| `add_channel` | yes | Appends channel ID to `discord.channel_ids` in config and saves |
| `remove_channel` | yes | Removes channel ID from `discord.channel_ids` in config and saves |

---

## EventHandle

`bot/handles/event_handle.py` — receives an `AdmineMessage` from the Pub/Sub subscription and dispatches based on the message's `tags` list. A single message can trigger multiple handlers.

| Tag | Handler | Discord notification |
|---|---|---|
| `server_on` | `__server_on` | `"Server has started with message: <payload>"` |
| `server_off` | `__server_off` | `"Server has stopped with message: <payload>"` |
| `notification` | `__notification` | The message payload verbatim |
| `new_server_ips` | `__new_server_ips` | `"Received new server IPs: <ip1,ip2>"` |
| `mod_install_result` | `__mod_install_result` | `"📦 Mod Install Result: <payload>"` |

---

## Configuration

You only need to define the `discord` section. All other sections fall back to defaults.

**Minimal config (required only):**

```json
{
    "discord": {
        "token": "YOUR_DISCORD_BOT_TOKEN",
        "commandprefix": "!mc",
        "administrators": ["admin_user_id"],
        "channel_ids": ["allowed_channel_id"]
    }
}
```

**Full reference with defaults:**

```json
{
    "logging": {
        "level": "INFO",
        "file": "/tmp/admine/logs/bot.log"
    },
    "security": {
        "ssl_verify": false
    },
    "providers": {
        "messaging": "DISCORD",
        "pubsub":    "REDIS",
        "minecraft": "REST",
        "vpn":       "REST"
    },
    "discord": {
        "token":          "YOUR_DISCORD_BOT_TOKEN",
        "commandprefix":  "!mc",
        "administrators": [],
        "channel_ids":    []
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

> **Note:** SSL certificate verification is disabled by default. Set `"security": {"ssl_verify": true}` to enable it.

Config is read from `./bot_config.json` by default. The path can be overridden as a constructor argument. Runtime changes (adding admins, channels) are persisted back to the same file via `Config.save()`.

---

## Build and test

```bash
make install      # Install dependencies via Poetry
make run          # Start the bot
make test         # Run pytest suite
make check        # Lint + format check (Ruff)
make fix          # Apply all auto-fixable Ruff fixes
make git-hooks    # Install pre-commit hooks
```

Run a single test:

```bash
poetry run pytest tests/test_command_handle.py -v
```
