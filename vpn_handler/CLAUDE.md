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

The VPN handler is a Rust/Actix-Web service that manages VPN network membership. It supports **ZeroTier** and **Tailscale** as providers, switchable via `vpn_config.vpn_type` in the TOML config. It has two entry points that run concurrently:

1. **REST API** (`src/api/`) — Actix-Web server wired in `server.rs`, handlers in `services.rs`. Three endpoints:
   - `GET /server-ips` — Returns current server's VPN IPs (looks up stored member ID, queries VPN provider)
   - `POST /auth-member` — Authorizes a member by ID
   - `GET /vpn-id` — Returns the configured network ID / tailnet slug

2. **Queue handler** (`src/queue_handler.rs`) — Blocking loop that receives Redis Pub/Sub messages and reacts to them. Two message flows:
   - `origin="server"` + tag `server_on`: authorizes the new member, fetches its IPs with retry logic, publishes them to `vpn_channel`, and deletes the old server member from the VPN provider.
   - `origin="bot"` + tag `auth_member`: authorizes a member and publishes a success message back to `vpn_channel`.

### Trait-based abstraction layers

All external dependencies are hidden behind traits, enabling mockall-based unit testing without live services:

| Trait | Location | Implementations |
|---|---|---|
| `TVpnClient` | `src/vpn/vpn.rs` | `ZerotierVpn` (wraps `zerotier-central-api` crate), `TailscaleVpn` (reqwest) |
| `KeyValueStore` | `src/persistence/key_value_storage.rs` | `SledStore` (embedded Sled DB) |
| `TSubscriber` / `TPublisher` | `src/pub_sub/pub_sub.rs` | `RedisPubSub` |

Factories (`vpn_factory.rs`, `key_value_storage_factory.rs`, `pub_sub_factory.rs`) select implementations based on enum variants from config.

### Adding a new VPN provider

1. Create `src/vpn/<provider>_vpn.rs` implementing `TVpnClient`.
2. Add a variant to `VpnType` in `src/vpn/vpn_factory.rs`.
3. Add a match arm in `VpnFactory::create_vpn`.
4. Declare the module in `src/vpn/mod.rs`.
5. If the provider needs a different default `api_url`, add normalization in `Config::new()` in `src/config.rs`.

No other files need to change.

### Config

Loaded from TOML via `src/config.rs`. All fields have defaults (see `impl Default for ...`). Required fields per provider:

- **ZeroTier**: `vpn_config.api_key` (Central API key) and `vpn_config.network_id` (network hex ID).
- **Tailscale**: `vpn_config.api_key` (`tskey-api-...`) and `vpn_config.network_id` (tailnet slug, e.g. `example.com`).

`api_url` defaults to the canonical base URL for each provider and only needs to be set if using a self-hosted or proxy endpoint.

Logger config is read from `./etc/log4rs.yaml`; if missing, it falls back to stdout + `/tmp/admine/logs/vpn_handler.log`.

### Persistence

Sled (embedded key-value store) persists `server_member_id` across restarts. This is used to detect and remove the previous server's VPN membership when a new `server_on` event arrives.

### Message format

All Redis messages are JSON-serialized `AdmineMessage` (`src/models/admine_message.rs`), which carries an `origin` string, a `tags` list, and a `message` payload. Routing in `queue_handler.rs` dispatches on `origin` first, then checks tags.

The `message` field in `server_on` and `auth_member` events carries the provider-specific member identifier:
- ZeroTier: 10-char hex node ID
- Tailscale: numeric device ID
