# Server Handler

The server_handler manages the full lifecycle of a Minecraft server running inside Docker. It:

- Subscribes to Redis Pub/Sub channels and reacts to lifecycle commands (start, stop, restart, etc.)
- Exposes a REST API for status queries, console access, and mod management
- Dynamically generates a `docker-compose.yaml` from a Go template on every `Start` — no manual Docker file editing required

---

## Package structure

```
server_handler/
├── cmd/server_handler/main.go   # Wiring only: constructs all objects and starts the app
└── internal/
    ├── server/                  # MinecraftServer interface + dockerMinecraftServer + domain models
    ├── pubsub/                  # PubSubService interface + Redis impl + EventHandler
    ├── api/
    │   ├── routes.go            # Gin router wiring
    │   ├── server.go            # HTTP server start/stop
    │   └── handlers/            # Stateless HTTP handlers (server.go, mod.go)
    ├── api/models/              # Request/response structs (no business logic)
    ├── deployment/              # docker-compose.yaml template (go:embed) + renderer
    ├── docker/                  # Docker SDK helpers: exec, log streaming, compose wrapper
    ├── logger/                  # slog setup (file + stdout, configurable level)
    ├── config/                  # Config struct + YAML loader + defaults
    └── testutils/               # Shared testify mocks (MockMinecraftServer, MockPubSubService)
```

All layers use explicit constructor injection. There is no global state.

---

## Dependency flow

```
main.go
  └─ config.LoadConfig()
  └─ server.NewDocker(cfg)          ─► dockerMinecraftServer
  └─ pubsub.NewRedis(cfg, ctx)      ─► redisPubSub
  └─ pubsub.NewEventHandler(...)    ─► EventHandler
  └─ api.SetupRouter(...)           ─► gin.Engine
  └─ api.StartServer(router)        ─► HTTP server
  └─ EventHandler.Listen()          ─► blocking pub/sub loop
```

---

## MinecraftServer interface

Defined in `internal/server/server.go`. Every operation that touches the server goes through this interface, making handlers and event listeners testable without Docker.

```go
type MinecraftServer interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Down(ctx context.Context) error
    Restart(ctx context.Context) error
    Status(ctx context.Context) (*ServerStatus, error)
    Info(ctx context.Context) (*ServerInfo, error)
    Logs(ctx context.Context, n int) ([]string, error)
    StartUpInfo(ctx context.Context) string
    ExecuteCommand(ctx context.Context, command string) (*CommandResult, error)
    InstallMod(ctx context.Context, fileName string, modData io.Reader) (*ModInstallResult, error)
    ListMods(ctx context.Context) (*ModListResult, error)
    RemoveMod(ctx context.Context, fileName string) (*ModInstallResult, error)
}
```

The only implementation is `dockerMinecraftServer` (`internal/server/docker.go`).

### Key implementation details

**Start** — renders the Compose template, then calls `docker compose up --detach`.

**Stop** — sends RCON `/stop`, streams container logs until `"All dimensions are saved"`, then calls `docker compose stop`. Falls back to `docker compose down` if RCON is unreachable.

**Restart** — calls `Stop` then `Start` with separate contexts (each phase has its own timeout, injected by the EventHandler).

**Info** — reads `minecraftVersion` and `modEngine` from config; queries `javaVersion` live via `docker compose exec java -version`; queries `maxPlayers` and `seed` live via RCON (`list` and `seed` commands).

**Status** — queries RCON `list` (online/offline check + player count), `forge tps` / `mspt` (TPS), and `stat /proc/1` inside the container (uptime).

---

## EventHandler (Pub/Sub)

`internal/pubsub/handler.go` routes incoming `AdmineMessage` tags to `MinecraftServer` methods.

| Tag | Action | Timeout used |
|---|---|---|
| `server_on` | `server.Start()` | `ServerOnTimeout` |
| `server_off` | `server.Stop()` | `ServerOffTimeout` |
| `server_down` | RCON `/stop` + `server.Down()` | `ServerCommandExecTimeout` + `ServerOffTimeout` |
| `restart` | `server.Stop()` then `server.Start()` | `ServerOffTimeout` + `ServerOnTimeout` |
| `command` | `server.ExecuteCommand()` | `ServerCommandExecTimeout` |

The handler ignores messages whose `origin` matches `app.self_origin_name`, preventing feedback loops.

After a successful `server_on` or `restart`, publishes a `server_on` message with the ZeroTier node ID (from `StartUpInfo`).

---

## REST API

Base URL: `http://<host>:<port>/api/v1`

> **Note:** server Start/Stop/Restart are **not** exposed via REST. They are triggered exclusively through Redis Pub/Sub.

| Method | Path | Description |
|---|---|---|
| `GET` | `/info` | Server info (version, java, mod engine, max players, seed) |
| `GET` | `/status` | Live status (health, TPS, uptime) |
| `GET` | `/logs?n=<int>` | Last N log lines (max 100, default 100) |
| `POST` | `/command` | Execute RCON command — body: `{"command":"<cmd>"}` |
| `GET` | `/resources` | Host CPU, memory, and disk usage |
| `POST` | `/mods` | Install mod — multipart `file` field or JSON `{"url":"..."}` |
| `GET` | `/mods` | List `.jar` files in `/data/mods/` |
| `DELETE` | `/mods/:filename` | Remove a mod by filename |
| `GET` | `/health` | Health check (no auth, always 200) |

### Response shapes

**`GET /info`**
```json
{
  "minecraftVersion": "1.20.1",
  "javaVersion": "17.0.9",
  "modEngine": "Fabric 0.15.11",
  "maxPlayers": 20,
  "seed": "1234567890"
}
```

**`GET /status`**
```json
{
  "health": "HEALTHY",
  "status": "ONLINE",
  "description": "Server is online - There are 2 of a max of 20 players online",
  "uptime": "2h 15m",
  "tps": 19.8
}
```

**`GET /resources`**
```json
{
  "cpu_usage": 12.4,
  "memory_used": 4294967296,
  "memory_total": 17179869184,
  "memory_used_percent": 25.0,
  "disk_used": 21474836480,
  "disk_total": 107374182400,
  "disk_used_percent": 20.0
}
```

**`POST /mods` and `DELETE /mods/:filename`**
```json
{ "file_name": "fabric-api-0.92.jar", "success": true, "message": "Mod installed successfully" }
```

---

## Docker Compose generation

The template is at `internal/deployment/docker-compose.yaml.tmpl`, embedded in the binary at build time (`//go:embed`). It is re-rendered on every `Start` and written to `docker.compose_output_path`. Do not edit the output file — changes will be overwritten.

To customise the Compose structure (extra networks, labels, resource limits), edit the template and rebuild.

### Package structure

```
server_handler/
├── cmd/server_handler/main.go   # Wiring: constructs all objects and starts the app
└── internal/
    ├── server/                  # MinecraftServer interface, Docker implementation, domain models
    ├── pubsub/                  # PubSubService interface, Redis implementation, EventHandler
    ├── api/
    │   ├── routes.go            # Gin router setup
    │   └── handlers/            # HTTP handlers (server.go, mod.go) — stateless, deps via constructor
    ├── deployment/              # docker-compose.yaml template rendering
    ├── docker/                  # Low-level Docker SDK helpers (exec, log tailing, compose)
    ├── logger/                  # slog setup
    └── config/                  # Config loading and defaults
```

All layers use explicit constructor injection — there is no global state.

---

## Configuration

Full reference with defaults:

```yaml
app:
  self_origin_name: "server"
  log_file_path: "/tmp/admine/logs/server_handler.log"
  log_level: "INFO"                  # DEBUG | INFO | WARN | ERROR

pubsub:
  type: "redis"
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
  admine_channels_map:
    server_channel:  "server_channel"
    command_channel: "command_channel"
    vpn_channel:     "vpn_channel"

minecraft_server:
  runtime_type: "docker"
  server_up_timeout:           "2m"
  server_off_timeout:          "1m"
  server_command_exec_timeout: "30s"
  mod_install_timeout:         "2m"
  rcon_address:  "127.0.0.1:25575"
  rcon_password: "admineRconPassword!"   # must match extra_env.RCON_PASSWORD

  docker:
    compose_output_path: "./generated/docker-compose.yaml"
    container_name: "mine_server"
    service_name:   "mine_server"
    data_path: "./minecraft-data"        # bind-mounted as /data inside the container

  image:

    # Server type – passed as the TYPE env var to itzg/docker-minecraft-server.
    # Common values: VANILLA, FABRIC, FORGE, NEOFORGE, PAPER, QUILT, PURPUR
    # Full list: https://docker-minecraft-server.readthedocs.io/en/latest/types-and-platforms/
    type: "FABRIC"

    # Minecraft version, e.g. "1.20.1". Use "LATEST" to always pull the newest release.
    version: "1.20.1"

    # Pin the Fabric loader version (leave empty for latest; only relevant when type: FABRIC).
    fabric_loader_version: ""
    forge_version: ""

    # Selects the JDK image tag (e.g. "java21", "java17"). Leave empty to use the itzg default.
    java_version: ""

    # URL to a modpack archive. When set, itzg downloads and installs it automatically.
    # Leave empty if you manage mods manually through the API.
    modpack_url: ""
    extra_env:
      RCON_PASSWORD: "admineRconPassword!"   # ← must match rcon_password above
      MEMORY: "4G"                           # JVM heap size
      # MAX_PLAYERS:       "20"
      # DIFFICULTY:        "normal"
      # MOTD:              "My Server"
      # ONLINE_MODE:       "false"
      # ENFORCE_WHITELIST: "true"
      # WHITELIST:         "Player1,Player2"
      # OPS:               "Player1"
      # VIEW_DISTANCE:     "10"
      # SPAWN_PROTECTION:  "0"

  zerotier:
    enabled: false
    network_id: ""
    container_name: "zerotier"
    api_secret: ""

web_server:
  host: "0.0.0.0"
  port: 3000
```

### Modpack platforms

**Modrinth:**
```yaml
image:
  type: "MODRINTH"

  version: "1.20.1"   # target Minecraft version, or LATEST
  extra_env:
    RCON_PASSWORD: "change-me!"
    MEMORY: "4G"
    MODRINTH_MODPACK: "fabric-api"          # slug, project ID, page URL, or .mrpack URL
    MODRINTH_VERSION: ""                    # specific version ID — omit for latest
    MODRINTH_LOADER: "fabric"              # fabric | forge | quilt — omit for auto-detect
    MODRINTH_MODPACK_VERSION_TYPE: "release" # release | beta | alpha
```

**CurseForge:**
```yaml
image:
  type: "AUTO_CURSEFORGE"
  extra_env:
    RCON_PASSWORD: "change-me!"
    MEMORY: "4G"          # CurseForge packs often need ≥4G
    CF_API_KEY: "your-curseforge-api-key"
    CF_SLUG: "all-the-mods-9"             # modpack slug from the CurseForge URL
    CF_FILE_ID: ""                         # pin a specific file ID — omit for latest
    # CF_PAGE_URL: "https://www.curseforge.com/minecraft/modpacks/all-the-mods-9"
```

**FTB:**
```yaml
image:
  type: "FTBA"
  extra_env:
    RCON_PASSWORD: "change-me!"
    MEMORY: "4G"
    FTB_MODPACK_ID: "31"          # numerical modpack ID
    FTB_MODPACK_VERSION_ID: ""    # specific version — omit for latest
```

---

## Build and test

```bash
make setup         # install golangci-lint + staticcheck
make build         # compile → ./bin/server_handler
make run           # build and run
make dev           # run with hot reload
make test          # run all tests
make test-coverage # generate coverage.html
make check         # vet + lint + staticcheck
make fmt           # go fmt
```

Run a single test or package:

```bash
go test -v ./internal/pubsub/
go test -v ./internal/api/handlers/ -run TestGetStatus_Success
```
