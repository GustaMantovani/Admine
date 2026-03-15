# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Admine** is an infrastructure management solution for Minecraft servers running on Linux. It orchestrates server lifecycle management, ZeroTier VPN connectivity, and a Discord-based administration interface.

The system is composed of four independently deployable components that communicate **exclusively through Redis Pub/Sub**:

- **`server_handler/`** (Go) — Orchestrates Minecraft server lifecycle via dynamically generated Docker Compose
- **`vpn_handler/`** (Rust) — Manages ZeroTier network membership and authorizations
- **`bot/`** (Python) — Discord bot; translates user commands to Pub/Sub messages
- **Redis** — The sole inter-component communication bus (channels: `server_channel`, `command_channel`, `vpn_channel`)

## Commands

### Server Handler (Go)

```bash
cd server_handler
make build          # Compile to ./bin/server_handler
make run            # Build and run
make dev            # Run with hot reload
make test           # Run tests
make test-coverage  # HTML coverage report
make check          # vet + lint + staticcheck
make fmt            # go fmt
make lint           # golangci-lint
make setup          # Install golangci-lint and staticcheck
```

### VPN Handler (Rust)

```bash
cd vpn_handler
cargo build         # Debug build
cargo build --release
cargo test
cargo fmt
```

### Discord Bot (Python — Poetry)

```bash
cd bot
make install        # Install via Poetry
make run            # Run bot
make test           # pytest tests/
make check          # lint + format-check (Ruff)
make lint           # Ruff check
make format         # Ruff format
make fix            # Apply all auto-fixable fixes
make git-hooks      # Install pre-commit hooks
```

### Redis (local dev)

```bash
cd pubsub/redis
docker compose up -d
```

## Architecture

### Inter-Component Communication

All coordination flows through Redis Pub/Sub. No component calls another directly. The bot publishes commands; handlers subscribe and respond back through Redis.

### Server Handler internals

- `internal/server/` — `MinecraftServer` interface + Docker implementation (RCON, lifecycle, mods, TPS). All domain models live here too (`ServerStatus`, `ServerInfo`, `CommandResult`, etc.)
- `internal/pubsub/` — `PubSubService` interface + Redis implementation + `EventHandler` (routes pubsub commands to server operations)
- `internal/api/` — Gin REST API: `routes.go` wires handlers; `handlers/server.go` and `handlers/mod.go` are stateless and receive all deps via constructor
- `internal/deployment/` — Dynamically generates `docker-compose.yaml` from a Go template on each Start (output goes to `generated/`, which is gitignored)
- `internal/docker/` — Low-level Docker SDK helpers (container exec, log tailing, compose wrapper)
- `internal/logger/` — slog setup (file + level)
- `cmd/server_handler/main.go` — Wiring only: constructs all objects, injects dependencies, starts web server and pubsub listener

All dependency injection is explicit via constructors — no global state.

### VPN Handler internals

- `src/vpn/` — ZeroTier Central API client
- `src/pub_sub/` — Redis integration
- `src/persistence/` — Sled embedded DB for local state
- `src/queue_handler.rs` — Event queue processing

### Bot internals

- `src/bot/handles/` — Discord command and event handlers
- `src/bot/external/` — Provider/abstraction pattern for services (Redis, server handler API, VPN handler API)
- `src/bot/models/` — Pydantic data models

## Configuration

| Component | File |
|-----------|------|
| Server Handler | `server_handler/server_handler_config.yaml` |
| VPN Handler | `vpn_handler/etc/vpn_handler_config.toml` |
| Bot | `bot/bot_config.json` |
| Redis | `pubsub/redis/redis.conf` |

## Git Conventions

Commits use [Gitmoji](https://gitmoji.dev/) prefixes (e.g., `:bug:` / `🐛` for fixes, `🔨` for refactors).

## Release

Built from repo root using Nushell:

```bash
nu utils/releasing/make-release.nu <version> [--clean] [--force] [--dev] [--push_tags]
```

Produces a deployment pack runnable via `./admine.sh start`.
