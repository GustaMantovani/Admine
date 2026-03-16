# Admine Deploy Pack — Configuration Guide

This guide covers everything you need to do before running `./admine.sh start`. All credentials stay in the three config files inside this directory — nothing else needs to change.

---

## 1. Discord Bot Configuration

Edit `bot/bot_config.json`:

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

**How to get these values:**

1. Go to [discord.com/developers/applications](https://discord.com/developers/applications), create an application, then go to the **Bot** tab and copy the token.
2. Enable the **Message Content**, **Server Members**, and **Presence** privileged gateway intents on the same page.
3. Enable Developer Mode in Discord (Settings → Advanced) and right-click any user or channel to copy their ID.

**SSL verification** is disabled by default (avoids CA issues in minimal environments). To enable it:

```json
{
    "security": { "ssl_verify": true },
    "discord": { "...": "..." }
}
```

---

## 2. VPN Configuration

Edit `vpn_handler/etc/vpn_handler_config.toml`. Choose **one** provider.

### Option A — ZeroTier

```toml
[vpn_config]
vpn_type   = "Zerotier"
api_key    = "your_zerotier_api_token"
network_id = "your_zerotier_network_id"
```

1. Create an account at [my.zerotier.com](https://my.zerotier.com).
2. Under **Account**, generate an API Access Token.
3. Create a network and copy its Network ID.

Also enable ZeroTier in `server_handler/server_handler_config.yaml`:

```yaml
minecraft_server:
  zerotier:
    enabled: true
    network_id: "your_zerotier_network_id"
```

### Option B — Tailscale

```toml
[vpn_config]
vpn_type   = "Tailscale"
api_key    = "tskey-api-..."
network_id = "your-tailnet-slug.ts.net"
```

1. Go to [login.tailscale.com/admin/settings/keys](https://login.tailscale.com/admin/settings/keys) and generate an **API access token** (for `api_key`).
2. Your tailnet slug is shown on the Settings page (e.g. `example.ts.net`).

Also enable Tailscale in `server_handler/server_handler_config.yaml`:

```yaml
minecraft_server:
  tailscale:
    enabled: true
    auth_key: "tskey-auth-..."     # reusable, ephemeral recommended
    hostname:  "minecraft-server"  # optional
```

Generate the `auth_key` at the same keys page — use type **Auth key**, enable **Reusable** and **Ephemeral**.

---

## 3. Minecraft Server Configuration

Edit `server_handler/server_handler_config.yaml`. At minimum set the server type, version, and RCON password:

```yaml
minecraft_server:
  rcon_password: "your_secure_rcon_password"
  image:
    type:    "FABRIC"    # VANILLA | FABRIC | FORGE | NEOFORGE | PAPER | MODRINTH | …
    version: "1.20.1"
    extra_env:
      RCON_PASSWORD: "your_secure_rcon_password"
      MEMORY: "4G"
```

The `rcon_password` and `RCON_PASSWORD` env var **must match exactly**. All other fields have sensible defaults.

The Minecraft server runs via [`itzg/docker-minecraft-server`](https://github.com/itzg/docker-minecraft-server). See its documentation for the full list of supported types and environment variables.

---

## 4. Start

```bash
./admine.sh start
```

This starts Redis (Docker Compose), then `server_handler`, `vpn_handler`, and `bot` as background processes. Logs go to `/tmp/admine/logs/`.

```bash
./admine.sh status                   # Show status of all services
./admine.sh logs                     # Last 50 lines from all logs
./admine.sh logs bot -f              # Follow bot logs
./admine.sh logs server_handler -n 100
./admine.sh restart vpn_handler      # Restart a single service
./admine.sh stop                     # Stop everything
```

---

## Discord Commands Reference

**Server Control**: `/on` `/off` `/restart`

**Monitoring**: `/status` `/info` `/resources` `/logs`

**Administration**: `/command <cmd>` `/adm <user>` `/add_channel` `/remove_channel`

**Mod Management**: `/install_mod` `/list_mods` `/remove_mod <file>`

**VPN**: `/auth <id>` `/vpn_id` `/server_ips`

---

For the full project documentation visit **[github.com/GustaMantovani/Admine](https://github.com/GustaMantovani/Admine)**.
