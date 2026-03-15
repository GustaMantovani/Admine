# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
cargo build               # Debug build
cargo build --release     # Release build
cargo test                # Run all tests
cargo test <test_name>    # Run a single test by name (substring match)
cargo fmt                 # Format code
cargo clippy              # Lint
```

The config file is read from `./etc/vpn_handler_config.toml` by default, or from a path given as the first CLI argument.

## Architecture

The VPN handler is a Rust/Actix-Web service that manages ZeroTier network membership. It has two entry points that run concurrently:

1. **REST API** (`src/api/`) — Actix-Web server wired in `server.rs`, handlers in `services.rs`. Three endpoints:
   - `GET /server-ips` — Returns current server's VPN IPs (looks up stored member ID, queries ZeroTier)
   - `POST /auth-member` — Authorizes a member by ID
   - `GET /vpn-id` — Returns the configured ZeroTier network ID

2. **Queue handler** (`src/queue_handler.rs`) — Blocking loop that receives Redis Pub/Sub messages and reacts to them. Two message flows:
   - `origin="server"` + tag `server_on`: authorizes the new member, fetches its IPs with retry logic, publishes them to `vpn_channel`, and deletes the old server member from ZeroTier.
   - `origin="bot"` + tag `auth_member`: authorizes a member and publishes a success message back to `vpn_channel`.

### Trait-based abstraction layers

All external dependencies are hidden behind traits, enabling mockall-based unit testing without live services:

| Trait | Location | Implementations |
|---|---|---|
| `TVpnClient` | `src/vpn/vpn.rs` | `ZerotierVpn` (wraps `zerotier-central-api` crate) |
| `KeyValueStore` | `src/persistence/key_value_storage.rs` | `SledStore` (embedded Sled DB) |
| `TSubscriber` / `TPublisher` | `src/pub_sub/pub_sub.rs` | `RedisPubSub` |

Factories (`vpn_factory.rs`, `key_value_storage_factory.rs`, `pub_sub_factory.rs`) select implementations based on enum variants from config.

### Config

Loaded from TOML via `src/config.rs`. All fields have defaults (see `impl Default for ...`), so only `vpn_config.api_key` and `vpn_config.network_id` are strictly required for ZeroTier to function. Logger config is read from `./etc/log4rs.yaml`; if missing, it falls back to stdout + `/tmp/admine/logs/vpn_handler.log`.

### Persistence

Sled (embedded key-value store) persists `server_member_id` across restarts. This is used to detect and remove the previous server's ZeroTier membership when a new `server_on` event arrives.

### Message format

All Redis messages are JSON-serialized `AdmineMessage` (`src/models/admine_message.rs`), which carries an `origin` string, a `tags` list, and a `message` payload. Routing in `queue_handler.rs` dispatches on `origin` first, then checks tags.
