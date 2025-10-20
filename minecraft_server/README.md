# Minecraft Server with ZeroTier

 Docker-based Minecraft server solution with VPN connectivity provided by ZeroTier. Supports Forge, Fabric, and All The Mods (ATM) modpack.

## Available Server Types

### Forge Server
Located in `forge/` directory - supports Forge mods for enhanced gameplay.

### Fabric Server  
Located in `fabric/` directory - supports Fabric mods with better performance.

### ATM Server (All The Mods)
Located in `atm/` directory - packaged modpack with pre-configured mods and scripts.

## Quick Start

1. **Choose your server type** (forge, fabric, or atm)
2. **Navigate to the directory**:
   ```bash
   cd forge/  # or cd fabric/ or cd atm/
   ```

3. **Set up configuration files**:
   ```bash
   ./setup.sh
   ```

4. **Configure environment variables** in `.env`:
   ```bash
   NETWORK_ID=your_zerotier_network_id
   JAVA_VERSION=17
   MINECRAFT_VERSION=1.20.1
   # For Forge:
   FORGE_VERSION=1.20.1-47.4.0
   # For Fabric:
   FABRIC_VERSION=0.14.21
   FRABRIC_INSTALLER_VERSION=0.11.2
   # For ATM:
   ATM_SERVER_FILES_DOWNLOAD_URL=https://mediafilez.forgecdn.net/files/6986/129/Server-Files-1.1.0.zip
   SERVER_FILES_VERSION=1.1.0
   ```

5. **Start the server**:
   ```bash
   docker-compose up -d
   ```

## Configuration Files

### Required Files (created by setup.sh)
- `config/eula.txt` - Minecraft EULA agreement
- `config/server.properties` - Server configuration
- `config/user_jvm_args.txt` - JVM memory and performance settings

### Data Directories
- `data/world/` - World save files
- `data/player-management/` - Player data (ops, whitelist, bans)
- `data/cache/` - Server cache files
- `mods/` - Mod files (JAR format)

Note for ATM: the modpack includes its own set of mods and scripts; you typically don't need to add mods manually. Use `atm/`'s `docker-compose.yaml` and `Dockerfile`.

## Network Configuration

The server runs with ZeroTier VPN integration:
- **Port 25565**: Minecraft game port
- **Port 25575**: RCON port (for remote commands)
- **Port 9993**: ZeroTier VPN port

## Adding Mods

1. Download mod files (.jar format)
2. Place them in the `mods/` directory
3. Restart the server: `docker-compose restart`

## Server Management

### Start server
```bash
docker-compose up -d
```

### Stop server
```bash
docker-compose down
```

### View logs
```bash
docker-compose logs -f mine_server
```

### ATM specifics
- The `atm/` image downloads the modpack server files automatically using `ATM_SERVER_FILES_DOWNLOAD_URL` and `SERVER_FILES_VERSION`.
- If you mount `config/server.properties` via volume (as in `atm/docker-compose.yaml`), your settings are preserved by the entrypoint.

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `NETWORK_ID` | ZeroTier network ID | `1c33c1ced0613a58` |
| `JAVA_VERSION` | Java version | `17` |
| `MINECRAFT_VERSION` | Minecraft version | `1.20.1` |
| `FORGE_VERSION` | Forge version (forge only) | `1.20.1-47.4.0` |
| `FABRIC_VERSION` | Fabric loader version (fabric only) | `0.14.21` |
| `FRABRIC_INSTALLER_VERSION` | Fabric installer version (fabric only) | `0.11.2` |
| `ATM_SERVER_FILES_DOWNLOAD_URL` | ATM server files zip URL (atm only) | `https://mediafilez.forgecdn.net/files/6986/129/Server-Files-1.1.0.zip` |
| `SERVER_FILES_VERSION` | ATM server files version (atm only) | `1.1.0` |
