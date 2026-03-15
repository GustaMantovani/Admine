# Admine Developer Guide

## Overview

Admine is a distributed system with four independently deployable components. Each component has its own language, toolchain, and configuration. They communicate exclusively through Redis Pub/Sub — no direct calls between services.

---

## Architecture

![Admine](.readme/Admine.png)

### Components

#### 1. Server Handler (Go)
Orchestrates the Minecraft server lifecycle through Docker Compose. On every `Start`, it renders `internal/deployment/docker-compose.yaml.tmpl` with values from `server_handler_config.yaml` and calls `docker compose up`. Supports all [`itzg/docker-minecraft-server`](https://github.com/itzg/docker-minecraft-server) server types (Vanilla, Fabric, Forge, NeoForge, Paper, Modrinth, CurseForge, FTB) and an optional ZeroTier sidecar container. Exposes a REST API consumed by the bot.

#### 2. VPN Handler (Rust)
Manages ZeroTier network membership via the ZeroTier Central REST API. Handles member authorization and network queries. Uses a local Sled embedded database for persistence and an async Redis PubSub client for event handling.

#### 3. Discord Bot (Python)
The user-facing interface. Translates Discord slash commands into Pub/Sub messages and relays responses back to Discord channels. Built with discord.py.

#### 4. Message Bus (Redis)
All inter-component communication flows through three Redis Pub/Sub channels:

| Channel | Purpose |
|---|---|
| `server_channel` | Server lifecycle events (up, down, restarting) |
| `command_channel` | Command routing and results |
| `vpn_channel` | VPN state changes and member updates |

Message envelope (`AdmineMessage`):

```json
{
  "origin":  "component_name",
  "tags":    ["event_tag"],
  "message": "payload"
}
```

A component ignores messages whose `origin` matches its own name, preventing feedback loops.

---

## Repository structure

```
Admine/
├── server_handler/              # Go — lifecycle management + REST API
│   ├── cmd/server_handler/      # Binary entrypoint (wiring only)
│   └── internal/
│       ├── api/                 # Gin router, handlers (server, mod), response models
│       ├── config/              # YAML config struct + loader
│       ├── deployment/          # Docker Compose template rendering (go:embed)
│       ├── docker/              # Docker SDK helpers and compose wrapper
│       ├── logger/              # slog setup (file + level)
│       ├── pubsub/              # PubSubService interface, Redis impl, EventHandler
│       ├── server/              # MinecraftServer interface, dockerMinecraftServer, domain models
│       └── testutils/           # Shared testify mocks
├── vpn_handler/                 # Rust — ZeroTier integration
│   ├── src/
│   │   ├── api/                 # Actix-web server and service handlers
│   │   ├── config.rs            # TOML config struct
│   │   ├── models/              # Request/response types and AdmineMessage
│   │   ├── persistence/         # Sled key-value store abstraction + factory
│   │   ├── pub_sub/             # Async Redis PubSub (subscriber + publisher)
│   │   ├── vpn/                 # ZeroTier Central API client + factory
│   │   ├── queue_handler.rs     # Event queue: routes pub/sub messages to handlers
│   │   └── main.rs              # Wiring: constructs all deps, starts server + queue
│   └── etc/                     # vpn_handler_config.toml, log4rs.yaml
├── bot/                         # Python — Discord bot
│   ├── src/
│   │   ├── main.py              # Entrypoint: config, logging, signal handling
│   │   └── bot/
│   │       ├── bot.py           # Bot lifecycle (start, shutdown)
│   │       ├── config.py        # JSON config loader with defaults
│   │       ├── handles/         # command_handle.py, event_handle.py
│   │       ├── models/          # Data models (ServerInfo, Status, ResourceUsage, etc.)
│   │       └── services/        # Service layer grouped by domain:
│   │           ├── messaging/   # Discord client (slash commands, response formatters)
│   │           ├── minecraft/   # server_handler REST API client
│   │           ├── pubsub/      # Redis pub/sub client
│   │           └── vpn/         # vpn_handler REST API client
│   └── tests/                   # pytest test suite
├── pubsub/redis/                # Redis config and compose file for local dev
└── utils/
    ├── releasing/               # Release scripts (make-release.nu)
    │   └── templates/           # Admine-Deploy-Pack deployment template
    ├── pubsub/                  # Debugging scripts for pub/sub messages
    └── mocks/apis/              # Mock API compose files for local dev
```

---

## Local development

### Server Handler (Go)

```bash
cd server_handler
make setup       # install golangci-lint and staticcheck
make build       # compile → ./bin/server_handler
make run         # build and run
make dev         # run with hot reload
make test        # run tests
make check       # vet + lint + staticcheck
make fmt         # go fmt
```

Requires: Go 1.21+, Docker with Compose plugin.

### VPN Handler (Rust)

```bash
cd vpn_handler
cargo build              # debug build
cargo build --release    # release build
cargo test               # run tests
cargo fmt                # format
./target/debug/vpn_handler
```

Requires: Rust toolchain (stable), a running Redis instance.

### Bot (Python)

```bash
cd bot
make install      # install dependencies via Poetry
make test         # run pytest suite
make run          # start the bot
make check        # lint + format check (Ruff)
make fix          # apply all auto-fixable fixes
make git-hooks    # install pre-commit hooks
```

Requires: Python 3.11+, Poetry.

### Redis (local)

```bash
cd pubsub/redis
docker compose up -d
```

---

## Testing

### Server Handler

Tests use `testutils.MockMinecraftServer` and `testutils.MockPubSubService` (both built with `testify/mock`). HTTP handler tests create a `gin.TestMode` context directly — no running server needed. Pub/sub handler tests likewise inject mocks directly.

```bash
go test -v ./internal/pubsub/
go test -v ./internal/api/handlers/ -run TestGetStatus_Success
make test-coverage   # generates coverage.html
```

### Bot

Tests use `unittest.mock`.

```bash
cd bot && make test
```

---

## Commit conventions

This project uses [gitmoji](https://gitmoji.dev/) prefixes.

| Emoji | When to use |
|---|---|
| 🐛 | Bug fix |
| ✨ | New feature |
| 🔨 | Refactor |
| 📝 | Documentation |
| ✅ | Tests |
| 🔧 | Config / tooling |
| 🚀 | Deploy / release |

Example:
```
🐛 fix restart handler: call Start() after Stop() with correct timeouts
```

---

## Release

Releases are built with [Nushell](https://www.nushell.sh/) from the repo root:

```bash
nu utils/releasing/make-release.nu <version>

# Options
--clean        # run clean before each build
--force        # overwrite existing output and tags
--dev          # skip git tagging (local iterations)
--push_tags    # push annotated tag to origin (default: false)
--no_archive   # skip tar.gz/zip creation
```

The output is a self-contained `admine-deploy-pack-<os>-<arch>-<version>/` directory ready to be dropped on the target host and started with `./admine.sh start`.
