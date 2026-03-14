# Server Handler

The `server_handler` manages a Minecraft server running inside Docker. It:

- Listens to Redis Pub/Sub channels for server lifecycle commands (start, stop, restart, etc.)
- Exposes a REST API for server control, mod management, and status queries
- Dynamically generates a `docker-compose.yaml` using the [`itzg/docker-minecraft-server`](https://github.com/itzg/docker-minecraft-server) image based on centralized configuration — no manual editing of Docker files required

---

## Architecture overview

```
server_handler_config.yaml
         │
         ▼
  server_handler (Go)
         │
         ├─ generates ──► docker-compose.yaml  (from docker-compose.yaml.tmpl)
         │
         ├─ manages ───► itzg/minecraft-server container
         │
         └─ manages ───► zerotier/zerotier sidecar (optional)
```

On every `Start`, the handler renders the Go template
`internal/deployment/docker-compose.yaml.tmpl` with values from the config and writes the result to `docker.compose_output_path`. Docker Compose is then called against that file. This means **the config file is the single source of truth** — editing `docker-compose.yaml` by hand has no permanent effect.

---

## Configuration

All configuration lives in `server_handler_config.yaml`. Every field has a default value; you only need to override what differs from the defaults.

```yaml
app:
  self_origin_name: "server"
  log_file_path: "/tmp/admine/logs/server_handler.log"
  log_level: "INFO"            # DEBUG | INFO | WARN | ERROR

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
  server_up_timeout:            "2m"
  server_off_timeout:           "1m"
  server_command_exec_timeout:  "30s"
  mod_install_timeout:          "2m"

  # RCON – must match the value in image.extra_env.RCON_PASSWORD
  rcon_address:  "127.0.0.1:25575"
  rcon_password: "admineRconPassword!"   # ← change this

  docker:
    # Where the generated docker-compose.yaml is written (created on every Start).
    compose_output_path: "./generated/docker-compose.yaml"
    container_name: "mine_server"
    service_name:   "mine_server"
    # Host directory bind-mounted as /data inside the container.
    # The entire Minecraft working directory (world, mods, configs, …) is persisted here.
    data_path: "./minecraft-data"

  image:
    # Server type – passed as the TYPE env var to itzg/docker-minecraft-server.
    # Common values: VANILLA, FABRIC, FORGE, NEOFORGE, PAPER, QUILT, PURPUR
    # Full list: https://docker-minecraft-server.readthedocs.io/en/latest/types-and-platforms/
    type: "FABRIC"

    # Minecraft version, e.g. "1.20.1". Use "LATEST" to always pull the newest release.
    version: "1.20.1"

    # JVM heap size passed as the MEMORY env var.
    memory: "4G"

    # Pin the Fabric loader version (leave empty for latest; only relevant when type: FABRIC).
    fabric_loader_version: ""

    # Pin the Forge version (leave empty for latest; only relevant when type: FORGE).
    forge_version: ""

    # URL to a modpack archive. When set, itzg downloads and installs it automatically.
    # Leave empty if you manage mods manually through the API.
    modpack_url: ""

    # Arbitrary itzg environment variables forwarded verbatim to the container.
    # Any itzg feature not covered by the explicit fields above can be set here.
    # Full reference: https://docker-minecraft-server.readthedocs.io/en/latest/
    extra_env:
      RCON_PASSWORD: "admineRconPassword!"   # ← must match rcon_password above
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
    # Set to true to start a zerotier/zerotier sidecar alongside the Minecraft server.
    # The sidecar runs with network_mode: host, making the server reachable via the ZT IP.
    enabled: false
    # ZeroTier network ID to join (required when enabled: true).
    network_id: ""
    # Name assigned to the ZeroTier container.
    container_name: "zerotier"
    # Optional: passed as ZEROTIER_API_SECRET inside the sidecar.
    api_secret: ""

web_server:
  host: "0.0.0.0"
  port: 3000
```

### Modloader / modpack selection

| Scenario | Fields to set |
|---|---|
| Vanilla | `type: VANILLA` |
| Fabric (latest loader) | `type: FABRIC` |
| Fabric (pinned loader) | `type: FABRIC`, `fabric_loader_version: "0.15.11"` |
| Forge (latest) | `type: FORGE` |
| Forge (pinned) | `type: FORGE`, `forge_version: "47.2.0"` |
| NeoForge | `type: NEOFORGE` |
| Paper | `type: PAPER` |
| Modpack via URL | `modpack_url: "https://…/pack.zip"` (type must match the pack's loader) |
| Modrinth modpack | see below |
| CurseForge modpack | see below |
| FTB modpack | see below |
| Extra itzg env vars | `extra_env: { KEY: value }` |

### Modpack platforms

Modpack platforms are configured entirely through `extra_env`. Each platform sets `TYPE` to a specific value and uses its own set of environment variables.

#### Modrinth

```yaml
image:
  type: "MODRINTH"
  version: "1.20.1"   # target Minecraft version, or LATEST
  memory: "4G"
  extra_env:
    RCON_PASSWORD: "change-me!"
    MODRINTH_MODPACK: "fabric-api"          # slug, project ID, page URL, or .mrpack URL
    MODRINTH_VERSION: ""                    # specific version ID — omit for latest
    MODRINTH_LOADER: "fabric"              # fabric | forge | quilt — omit for auto-detect
    MODRINTH_MODPACK_VERSION_TYPE: "release" # release | beta | alpha
```

#### CurseForge

Requires a free API key from [console.curseforge.com](https://console.curseforge.com/).

```yaml
image:
  type: "AUTO_CURSEFORGE"
  memory: "4G"          # CurseForge packs often need ≥4G
  extra_env:
    RCON_PASSWORD: "change-me!"
    CF_API_KEY: "your-curseforge-api-key"
    CF_SLUG: "all-the-mods-9"             # modpack slug from the CurseForge URL
    CF_FILE_ID: ""                         # pin a specific file ID — omit for latest
    # CF_PAGE_URL: "https://www.curseforge.com/minecraft/modpacks/all-the-mods-9"
```

#### Feed the Beast (FTB)

```yaml
image:
  type: "FTBA"
  memory: "4G"
  extra_env:
    RCON_PASSWORD: "change-me!"
    FTB_MODPACK_ID: "31"          # numerical modpack ID
    FTB_MODPACK_VERSION_ID: ""    # specific version — omit for latest
```

### Volume persistence

The host directory `docker.data_path` is bind-mounted to `/data` inside the container. The itzg image keeps the entire Minecraft working directory under `/data` (world, `server.properties`, mods, plugins, configs, logs, etc.), so everything is automatically persisted across restarts.

### ZeroTier sidecar

When `zerotier.enabled: true`, a `zerotier/zerotier` container is added to the generated Compose file. It runs with `network_mode: host` and `NET_ADMIN` / `SYS_ADMIN` capabilities so it can create the `tun` device and attach to the host network stack. The Minecraft container declares a `depends_on` on the ZeroTier service, ensuring ZeroTier is up before the server starts.

---

## Docker Compose generation

The template lives at:

```
server_handler/internal/deployment/docker-compose.yaml.tmpl
```

It is embedded in the binary at build time via `//go:embed`. When you need to customise the generated Compose structure (e.g. add extra networks, labels, or resource limits) edit that template file and rebuild — no Go code changes required.

The output file (`docker.compose_output_path`) is re-generated on every `Start`. Do **not** edit it manually; changes will be overwritten.

---

## REST API

Base path: `http://<host>:<port>/api/v1`

| Method | Path | Description |
|---|---|---|
| `POST` | `/server/start` | Start the server (generates Compose file, then `docker compose up`) |
| `POST` | `/server/stop` | Gracefully stop the server (`/stop` → waits for save → `docker compose stop`) |
| `POST` | `/server/restart` | Stop then start |
| `DELETE` | `/server` | Tear down containers (`docker compose down`) |
| `GET` | `/server/status` | Online/offline status, TPS, uptime |
| `GET` | `/server/info` | Static server info sourced from config (version, type, mod engine, max players) |
| `GET` | `/server/logs?n=<lines>` | Last N log lines from the container |
| `POST` | `/server/command` | Execute an RCON command; body: `{"command":"<cmd>"}` |
| `POST` | `/mods` | Upload a mod `.jar` (multipart); installs to `/data/mods/` |
| `GET` | `/mods` | List installed `.jar` files in `/data/mods/` |
| `DELETE` | `/mods/:filename` | Remove a mod by filename |

### `/server/info` response

Server info is read entirely from the configuration file — no container queries are made. Fields returned:

| Field | Source |
|---|---|
| `minecraftVersion` | `minecraft_server.image.version` |
| `modEngine` | Derived from `image.type` + `fabric_loader_version` / `forge_version` |
| `maxPlayers` | `image.extra_env.MAX_PLAYERS` (falls back to `-1` if not set) |
| `seed` | Always `"N/A - Seed Hidden"` (not tracked in config) |
| `javaVersion` | Always `"N/A - Not tracked in config"` |

### `/server/status` response

Status is determined at query time via RCON and Docker:

| Field | Source |
|---|---|
| `health` | RCON reachability |
| `status` | Online / Offline |
| `message` | Output of the `list` command |
| `uptime` | `stat /proc/1` inside the container |
| `tps` | `forge tps` or `mspt` RCON commands; defaults to `20.0` |

---

## Redis Pub/Sub

The handler subscribes to the channels defined under `pubsub.admine_channels_map`. Messages follow the Admine envelope format and trigger the same lifecycle operations as the REST API. The `self_origin_name` field is used to filter out messages that originated from this handler itself.

---

## Logs

Log files are written to `app.log_file_path` (default `/tmp/admine/logs/server_handler.log`). Set `app.log_level` to `DEBUG` to include RCON responses and container exec output.
