# Admine - Infrastructure Manager for Minecraft Servers

Admine is a comprehensive infrastructure management solution for Minecraft servers running on Linux systems. It provides automated server lifecycle management, VPN connectivity, and Discord-based administration interface.

## What is Admine?

Admine automates the complete management of Minecraft servers in a containerized environment, solving the common problem of running servers behind NAT (Network Address Translation) or in private networks where direct public access isn't available. The system is designed for server administrators who want to host Minecraft servers in private networks (home networks, cloud instances behind NAT, etc.) while providing secure and easy access to players through VPN connectivity.

The Minecraft server runs via the [`itzg/docker-minecraft-server`](https://github.com/itzg/docker-minecraft-server) Docker image.

## System Architecture

![Admine](.readme/Admine.png)

## Command Interface

Users can interact with the system through Discord commands:

**Server Control**
- `/on` - Start the server
- `/off` - Stop the server
- `/restart` - Restart the server

**Monitoring**
- `/status` - Get server status, TPS and uptime
- `/info` - Get server information (version, mod engine, player capacity)
- `/resources` - Get host resource usage (CPU, memory, disk)
- `/logs` - View the latest server logs

**Administration**
- `/command <cmd>` - Execute a Minecraft console command
- `/adm <user>` - Grant administrator permission to a Discord user
- `/add_channel` - Authorize a channel for bot interactions
- `/remove_channel` - Remove a channel authorization

**Mod Management**
- `/install_mod` - Install a mod from a file upload or URL
- `/list_mods` - List all installed mods
- `/remove_mod <file>` - Remove an installed mod

**VPN**
- `/auth <id>` - Authorize a VPN member
- `/vpn_id` - Get the VPN network ID
- `/server_ips` - Get the server's VPN IP addresses

---

## Table of Contents

- [Installation](#installation)
  - [1. Download the deploy pack](#1-download-the-deploy-pack)
  - [2. Configure Discord](#2-configure-discord)
  - [3. Configure VPN](#3-configure-vpn)
  - [4. Configure the Minecraft server](#4-configure-the-minecraft-server)
  - [5. Start](#5-start)
- [Running on free cloud environments](#running-on-free-cloud-environments)
  - [Oracle Cloud Free Tier](#oracle-cloud-free-tier)
  - [Google Cloud Shell](#google-cloud-shell)
  - [GitHub Codespaces](#github-codespaces)
  - [AWS CloudShell](#aws-cloudshell)
  - [Azure Cloud Shell](#azure-cloud-shell)
- [Mod logistics with Google Drive](#mod-logistics-with-google-drive)
- [Technical Documentation](#technical-documentation)
  - [Architecture overview](#architecture-overview)
  - [Message envelope](#message-envelope)
  - [Pub/Sub channels](#pubsub-channels)
  - [Message flows](#message-flows)
  - [Component documentation](#component-documentation)
- [Contributing](CONTRIBUTING.md)

---

## Installation

### 1. Download the deploy pack

Download the latest release from [GitHub Releases](https://github.com/GustaMantovani/Admine/releases). Each release ships a self-contained `admine-deploy-pack-<os>-<arch>-<version>` directory that includes compiled binaries for all components, config templates, and the `admine.sh` control script. No build step needed.

```bash
wget https://github.com/GustaMantovani/Admine/releases/download/v<version>/admine-deploy-pack-linux-x86_64-<version>.zip
unzip admine-deploy-pack-linux-x86_64-<version>.zip
cd admine-deploy-pack-linux-x86_64-<version>
```

The directory layout after extraction:

```
admine-deploy-pack/
├── admine.sh                    # Control script (start/stop/restart/status/logs)
├── bot/
│   ├── bot                      # Python bot binary
│   └── bot_config.json          # Bot configuration
├── server_handler/
│   ├── server_handler           # Go binary
│   └── server_handler_config.yaml
├── vpn_handler/
│   ├── vpn_handler              # Rust binary
│   └── etc/
│       └── vpn_handler_config.toml
└── pubsub/redis/
    └── docker-compose.yaml      # Redis container
```

### 2. Configure Discord

Edit `bot/bot_config.json`. The only required section is `discord`:

```json
{
    "discord": {
        "token": "YOUR_DISCORD_BOT_TOKEN",
        "commandprefix": "!mc",
        "administrators": ["your_discord_user_id"],
        "channel_ids": ["your_channel_id"]
    }
}
```

To create a bot and get your token: [Discord Developer Portal](https://discord.com/developers/applications). Enable the **Message Content**, **Server Members**, and **Presence** privileged gateway intents. To get user/channel IDs, enable Developer Mode in Discord (Settings → Advanced) and right-click any user or channel.

### 3. Configure VPN

Edit `vpn_handler/etc/vpn_handler_config.toml`. Set `vpn_type` to either `"Zerotier"` or `"Tailscale"` and fill in the corresponding credentials.

**Option A — ZeroTier**

```toml
[vpn_config]
vpn_type   = "Zerotier"
api_key    = "your_zerotier_api_token"
network_id = "your_zerotier_network_id"
```

Create an account at [my.zerotier.com](https://my.zerotier.com), generate an API Access Token under Account, and create a network to get the Network ID.

Also enable the ZeroTier sidecar in `server_handler/server_handler_config.yaml`:

```yaml
minecraft_server:
  zerotier:
    enabled: true
    network_id: "your_zerotier_network_id"
```

**Option B — Tailscale**

```toml
[vpn_config]
vpn_type   = "Tailscale"
api_key    = "tskey-api-..."
network_id = "your-tailnet-slug.ts.net"
```

Generate an API access token at [login.tailscale.com/admin/settings/keys](https://login.tailscale.com/admin/settings/keys). The `network_id` is your tailnet slug (visible on the Settings page).

Also enable the Tailscale sidecar in `server_handler/server_handler_config.yaml`:

```yaml
minecraft_server:
  tailscale:
    enabled: true
    auth_key: "tskey-auth-..."
    hostname: "minecraft-server"   # optional
```

Generate a Tailscale auth key at [login.tailscale.com/admin/settings/keys](https://login.tailscale.com/admin/settings/keys) (type: reusable, ephemeral recommended).

### 4. Configure the Minecraft server

Edit `server_handler/server_handler_config.yaml`. At minimum, set the server type, version, and RCON password:

```yaml
minecraft_server:
  rcon_password: "your_secure_rcon_password"   # must match RCON_PASSWORD below
  image:
    type:    "FABRIC"     # VANILLA | FABRIC | FORGE | NEOFORGE | PAPER | MODRINTH | …
    version: "1.20.1"
    extra_env:
      RCON_PASSWORD: "your_secure_rcon_password"
      MEMORY: "4G"
```

All other fields have sensible defaults. For the full configuration reference and supported server types, see [server_handler/README.md](server_handler/README.md).

The Minecraft server is powered by `itzg/docker-minecraft-server`. Refer to its documentation for the full list of supported server types and environment variables:

- [itzg/docker-minecraft-server — GitHub](https://github.com/itzg/docker-minecraft-server)
- [itzg/docker-minecraft-server — Docs](https://docker-minecraft-server.readthedocs.io)

### 5. Start

```bash
./admine.sh start
```

The script starts Redis (via Docker Compose), then `server_handler`, `vpn_handler`, and `bot` as background processes.

```bash
./admine.sh status                  # Show status of all services
./admine.sh logs                    # Last 50 lines from all logs
./admine.sh logs bot -f             # Follow bot logs
./admine.sh logs server_handler -n 100
./admine.sh restart vpn_handler     # Restart a single service
./admine.sh stop                    # Stop everything
```

Logs are written to `/tmp/admine/logs/`.

### Prerequisites

- Linux (x86_64 or arm64)
- Docker with Compose plugin
- A Discord bot token
- A VPN account: either [ZeroTier](https://www.zerotier.com) (API key + network ID) or [Tailscale](https://tailscale.com) (API key + auth key + tailnet slug)

---

## Running on free cloud environments

The typical problem with hosting a Minecraft server is needing a machine with a public IP. Most home connections and all free cloud environments sit behind NAT — the machine has no public IP and cannot accept inbound connections from the internet.

Admine was designed specifically for this scenario. Instead of requiring port forwarding or a public IP, it uses a VPN overlay network ([ZeroTier](https://www.zerotier.com) or [Tailscale](https://tailscale.com)): the server joins the VPN and players connect to it through that overlay IP. The server never needs to be reachable from the public internet.

This means you can run a Minecraft server for free on any of the following platforms:

---

### Oracle Cloud Free Tier

The best free option. Oracle's Always Free tier includes up to **4 Ampere A1 vCPUs and 24 GB RAM** on ARM64 VMs — enough to comfortably run the Minecraft server alongside all Admine components. Instances are persistent (no automatic shutdown) and have a real static public IP that is never used for game traffic.

- Sign up at [cloud.oracle.com](https://www.oracle.com/cloud/free/)
- Create an Always Free Compute instance (Ubuntu or Oracle Linux)
- Install Docker, clone the deploy pack, configure and run

### Google Cloud Shell

Google Cloud Shell provides a free persistent Linux environment directly in the browser, backed by a Debian VM with ~5 GB of home directory storage. It has Docker available and works well for running the Admine management components (bot, vpn_handler, server_handler). The shell has limited RAM (~1.7 GB) so it may struggle with heavier modpacks, but works for vanilla or light Fabric/Paper servers.

- Access at [shell.cloud.google.com](https://shell.cloud.google.com)
- No sign-up beyond a Google account required
- The environment pauses after inactivity but the home directory persists

### GitHub Codespaces

Codespaces gives every GitHub account **60 free hours/month** of a containerized dev environment with Docker, up to 4 vCPUs and 8 GB RAM depending on the machine type selected. Suitable for running the full Admine stack during active sessions.

- Access at [github.com/codespaces](https://github.com/codespaces)
- Open the Admine repository directly in a Codespace
- Note: Codespaces pause when inactive, so they work best for on-demand server sessions

### AWS CloudShell

AWS CloudShell is a free, browser-based shell pre-authenticated with your AWS account. It has 1 GB persistent storage but limited compute. Useful for running the Admine management components if you already have an AWS account.

- Access from the AWS Console toolbar
- Docker is available but requires root or a workaround depending on the region

### Azure Cloud Shell

Azure Cloud Shell provides a free Linux shell with 5 GB persistent storage mounted from an Azure Files share. Docker is available. Like AWS CloudShell, it is better suited for running the management components than for a RAM-heavy Minecraft server.

- Access at [shell.azure.com](https://shell.azure.com)
- Requires an Azure account (free tier available)

---

### Why it works behind NAT

The key is that **all outbound connections, no inbound**:

- The Minecraft server container joins the VPN network (ZeroTier or Tailscale) by making an outbound connection to the provider's coordination servers — no port forwarding needed.
- The bot connects to Discord via a persistent WebSocket (outbound).
- Redis is local to the machine.
- Players connect to the server using its VPN overlay IP after being authorized via `/auth`.

The host machine never needs to be directly reachable from the internet.

---

## Mod logistics with Google Drive

Managing a large modpack across multiple players and a remote server is one of the less glamorous parts of running a Minecraft server. A practical approach is to keep all mods in a shared Google Drive folder: share the link with players so they can download the client-side mods, and use that same link on the server to sync the server-side mods.

The standard tool for this is [gdown](https://github.com/wkentaro/gdown), a Python CLI that downloads public Google Drive files and folders — something that `curl` and `wget` cannot do reliably. The problem is that the original gdown downloads folder contents **sequentially**, making it slow for large modpacks with dozens of files.

To solve this, we maintain a fork with **parallel download support**:

**[GustaMantovani/gdown](https://github.com/GustaMantovani/gdown)** — adds a `--workers` flag to `download_folder`, spawning multiple threads to download files concurrently.

### Installing the fork

```bash
pip install git+https://github.com/GustaMantovani/gdown.git
```

### Downloading a modpack folder

```bash
# Sequential (original behavior, default workers=1)
gdown --folder "https://drive.google.com/drive/folders/<folder-id>"

# Parallel — use N concurrent workers
gdown --folder --workers 8 "https://drive.google.com/drive/folders/<folder-id>"

# Auto — uses all available CPU cores
gdown --folder --workers auto "https://drive.google.com/drive/folders/<folder-id>"

# Resume an interrupted download
gdown --folder --workers 8 --continue "https://drive.google.com/drive/folders/<folder-id>"
```

### Typical workflow

1. Upload your `.jar` mod files to a Google Drive folder and set sharing to **"Anyone with the link"**.
2. Share the link with your players so they can download client-side mods.
3. On the server, use `gdown` to pull the same folder into the mods directory before starting the server with `/on`.

> **Note:** Google Drive may throttle or block access after many rapid requests from the same IP. If you hit this, export your browser cookies to `~/.cache/gdown/cookies.txt` as described in the [gdown troubleshooting guide](https://github.com/wkentaro/gdown#faq).

---

## Technical Documentation

### Architecture overview

```
Discord ──► bot ──► Redis Pub/Sub ──► server_handler ──► Docker (Minecraft + VPN sidecar)
                         │
                         └──────────► vpn_handler ──────► VPN API (ZeroTier or Tailscale)
```

All inter-component communication flows exclusively through Redis Pub/Sub. No component calls another directly. The bot is the only entry point for user-initiated actions.

### Message envelope

Every Redis message is a JSON-serialized `AdmineMessage`:

```json
{
  "origin":  "component_name",
  "tags":    ["event_tag"],
  "message": "payload"
}
```

Each component ignores messages whose `origin` matches its own name, preventing feedback loops.

### Pub/Sub channels

| Channel | Used by | Purpose |
|---|---|---|
| `server_channel` | bot (pub), server_handler (sub) | Server lifecycle commands |
| `command_channel` | bot (pub), server_handler (sub), vpn_handler (sub) | Console commands and VPN auth requests |
| `vpn_channel` | vpn_handler (pub), bot (sub), server_handler (pub) | VPN state changes and server IP updates |

### Message flows

**Start server**

```
bot               Redis (server_channel)    server_handler
 │── server_on ──────────────────────────────► Start()
 │                                              ├─ generate docker-compose.yaml
 │                                              └─ docker compose up
 │
 │                                          server_handler → Redis (vpn_channel)
 │                                          server_on (VPN node ID) ──────────────────►
 │
 │                Redis (vpn_channel)       vpn_handler
 │◄── new_server_ips ◄────────────────────── process_server_up
 │                                              ├─ auth_member (with retry)
 │                                              ├─ get_member_ips (with retry)
 │                                              ├─ publish new_server_ips
 │                                              └─ delete old member + persist new ID
```

**Stop server**

```
bot               Redis (server_channel)    server_handler
 │── server_off ──────────────────────────────► Stop()
                                                ├─ RCON /stop (falls back to compose down)
                                                ├─ stream logs until "All dimensions are saved"
                                                └─ docker compose stop
```

**Authorize VPN member**

```
bot               Redis (command_channel)    vpn_handler
 │── auth_member ──────────────────────────────► auth_member(id) via VPN API
                                                   └─ publish auth_member_success
```

### Component documentation

Each component has its own detailed README:

- [server_handler/README.md](server_handler/README.md) — Go service: Docker lifecycle, RCON, mod management, pub/sub routing, config reference
- [vpn_handler/README.md](vpn_handler/README.md) — Rust service: ZeroTier/Tailscale integration, persistence, pub/sub routing, config reference
- [bot/README.md](bot/README.md) — Python bot: command routing, event handling, service abstractions, config reference
- [CONTRIBUTING.md](CONTRIBUTING.md) — Development guide: local setup, testing, commit conventions, release process
