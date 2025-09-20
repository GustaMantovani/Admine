# Minecraft Server with ZeroTier

Docker-based Minecraft server solution with VPN connectivity provided by ZeroTier. Supports both Forge and Fabric mod loaders.

## Available Server Types

### Forge Server
Located in `forge/` directory - supports Forge mods for enhanced gameplay.

### Fabric Server  
Located in `fabric/` directory - supports Fabric mods with better performance.

## Quick Start

1. **Choose your server type** (forge or fabric)
2. **Navigate to the directory**:
   ```bash
   cd forge/  # or cd fabric/
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

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `NETWORK_ID` | ZeroTier network ID | `1c33c1ced0613a58` |
| `JAVA_VERSION` | Java version | `17` |
| `MINECRAFT_VERSION` | Minecraft version | `1.20.1` |
| `FORGE_VERSION` | Forge version (forge only) | `1.20.1-47.4.0` |
| `FABRIC_VERSION` | Fabric loader version (fabric only) | `0.14.21` |
| `FRABRIC_INSTALLER_VERSION` | Fabric installer version (fabric only) | `0.11.2` |
