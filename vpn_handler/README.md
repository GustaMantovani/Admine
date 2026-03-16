# VPN Handler

The vpn_handler is a Rust/Actix-Web service that manages VPN network membership. It:

- Subscribes to Redis Pub/Sub channels and reacts to server lifecycle and user authorization events
- Exposes a REST API for direct member authorization and network queries
- Persists the current server's VPN member ID across restarts using an embedded Sled database

Supports **ZeroTier** and **Tailscale** as VPN providers, switchable via config without code changes.

---

## Package structure

```
vpn_handler/
├── src/
│   ├── main.rs                          # Wiring only: constructs all objects and starts the app
│   ├── config.rs                        # TOML config struct + loader + defaults
│   ├── errors.rs                        # VpnError enum
│   ├── queue_handler.rs                 # Event loop: routes pub/sub messages to handlers
│   ├── api/
│   │   ├── server.rs                    # Actix-Web server creation and wiring
│   │   └── services.rs                  # HTTP handlers (server_ip, auth_member, vpn_id)
│   ├── models/
│   │   ├── admine_message.rs            # Shared message envelope
│   │   ├── auth_member_request.rs       # POST /auth-member body
│   │   ├── server_ip_response.rs        # GET /server-ips response
│   │   ├── vpn_id_response.rs           # GET /vpn-id response
│   │   └── error_response.rs            # Error response shape
│   ├── persistence/
│   │   ├── key_value_storage.rs         # KeyValueStore trait + DynKeyValueStore alias
│   │   ├── sled_store.rs                # Sled implementation
│   │   └── key_value_storage_factory.rs # Factory: selects impl from config
│   ├── pub_sub/
│   │   ├── pub_sub.rs                   # TSubscriber + TPublisher traits + DynPubSub alias
│   │   ├── redis_pubsub.rs              # Async Redis implementation
│   │   └── pub_sub_factory.rs           # Factory: selects impl from config
│   └── vpn/
│       ├── vpn.rs                       # TVpnClient trait + DynVpn alias
│       ├── zerotier_vpn.rs              # ZeroTier Central API implementation
│       ├── tailscale_vpn.rs             # Tailscale API implementation
│       ├── vpn_factory.rs               # Factory: selects impl from config
│       └── public_ip.rs                 # Helper: fetch public IP
└── etc/
    ├── vpn_handler_config.toml          # Runtime configuration
    └── log4rs.yaml                      # Logger configuration
```

All layers use explicit constructor injection. There is no global state.

---

## Dependency flow

```
main.rs
  └─ Config::new()                   ─► Arc<Config>
  └─ StoreFactory::create(...)       ─► Arc<DynKeyValueStore>  (Sled)
  └─ VpnFactory::create(...)         ─► Arc<DynVpn>            (ZeroTier | Tailscale)
  └─ PubSubFactory::create(...)      ─► DynPubSub              (Redis)
  └─ pub_sub.subscribe([channels])
  └─ server::create_server(...)      ─► (ActixServer, ServerHandle)
  └─ Handle::new(...).run()          ─► blocking pub/sub loop  (spawned task)
  └─ actix_server.await              ─► HTTP server (blocking)
```

Graceful shutdown is coordinated via a `tokio::sync::watch` channel: a signal handler sends `true`, which stops the Actix server and the queue handler.

---

## Trait abstractions

All external dependencies are hidden behind traits, enabling mockall-based unit tests without live services:

| Trait | Location | Implementations |
|---|---|---|
| `TVpnClient` | `src/vpn/vpn.rs` | `ZerotierVpn` (wraps `zerotier-central-api` crate), `TailscaleVpn` (reqwest) |
| `KeyValueStore` | `src/persistence/key_value_storage.rs` | `SledStore` (embedded Sled DB) |
| `TSubscriber` / `TPublisher` | `src/pub_sub/pub_sub.rs` | `RedisPubSub` (async, multiplexed) |

Type aliases (`DynVpn`, `DynKeyValueStore`, `DynPubSub`) are `Box<dyn Trait + Send + Sync>`, making them injectable as `Arc<Dyn*>` in Actix handlers via `web::Data`.

---

## Queue handler (Pub/Sub)

`src/queue_handler.rs` — `Handle::run()` is a `tokio::select!` loop that either processes incoming messages or stops on shutdown signal.

Message routing dispatches on `origin` first, then on `tags`:

| Origin | Tag | Action |
|---|---|---|
| `server` | `server_on` | `process_server_up`: authorize member, fetch IPs with retry, publish `new_server_ips` to `vpn_channel`, delete old server member from VPN, persist new member ID |
| `bot` | `auth_member` | `process_auth_member`: authorize member, publish `auth_member_success` to `vpn_channel` |

### `server_on` flow in detail

1. Calls `vpn_client.auth_member(member_id)` — retries up to `retry_config.attempts` times with `retry_config.delay` between attempts.
2. Calls `vpn_client.get_member_ips_in_vpn(member_id)` — retries if empty (IP assignment may lag depending on provider).
3. Publishes an `AdmineMessage { origin: self_origin, tags: ["new_server_ips"], message: "ip1,ip2" }` to `vpn_channel`.
4. Reads `server_member_id` from Sled; if it differs from the new ID, calls `vpn_client.delete_member(old_id)` to clean up the stale node.
5. Writes the new `server_member_id` to Sled.

---

## REST API

The vpn_handler exposes a REST API consumed by the bot. For the full endpoint reference and response shapes, see the API specification in the repository.

---

## Persistence

Sled (embedded key-value store) persists `server_member_id` across restarts. This is the VPN member/device ID of the most recently known server instance. On every `server_on` event the old ID is deleted from the VPN provider and the new ID is saved.

---

## Configuration

Full reference with defaults (`etc/vpn_handler_config.toml`):

```toml
self_origin_name = "vpn"

[api_config]
host = "localhost"
port = 9000

[pub_sub_config]
url = "redis://localhost:6379"
pub_sub_type = "Redis"

[vpn_config]
vpn_type   = "Zerotier"   # "Zerotier" | "Tailscale"
api_url    = ""           # optional — defaults per provider (see below)
api_key    = ""           # required
network_id = ""           # required — ZeroTier: network ID hex; Tailscale: tailnet slug (e.g. "example.com")

[db_config]
path = "./etc/sled/vpn_store.db"
store_type = "Sled"

[admine_channels_map]
server_channel  = "server_channel"
command_channel = "command_channel"
vpn_channel     = "vpn_channel"

[retry_config]
attempts = 5
delay    = { secs = 3, nanos = 0 }
```

### VPN provider configuration

**ZeroTier** (default):
```toml
[vpn_config]
vpn_type   = "Zerotier"
api_key    = "<ZeroTier Central API key>"
network_id = "<16-char hex network ID>"
# api_url defaults to "https://api.zerotier.com/api/v1"
```

The `member_id` sent via pub/sub must be the ZeroTier node ID (10-char hex).

**Tailscale**:
```toml
[vpn_config]
vpn_type   = "Tailscale"
api_key    = "tskey-api-..."
network_id = "example.com"   # tailnet slug
# api_url defaults to "https://api.tailscale.com"
```

The `member_id` sent via pub/sub must be the Tailscale device ID (numeric string).

The config file is read from `./etc/vpn_handler_config.toml` by default, or from the path given as the first CLI argument. Logger config is read from `./etc/log4rs.yaml`; if missing, it falls back to stdout + `/tmp/admine/logs/vpn_handler.log`.

---

## Build and test

```bash
cargo build               # Debug build
cargo build --release     # Release build
cargo test                # Run all tests
cargo test <name>         # Run a single test (substring match)
cargo fmt                 # Format code
cargo clippy              # Lint
```
