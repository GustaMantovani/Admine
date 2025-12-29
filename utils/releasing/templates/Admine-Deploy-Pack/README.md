# Admine-slim Configuration Tutorial

This tutorial details the necessary configurations to run the Admine-slim system. Since this is a slim version for delivery, many things are already configured - you only need to add your personal credentials.

## 1. Discord Bot Configuration

In the `bot_config.json` file, you need to configure:

### Discord Bot Token:

1. Access the [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to the "Bot" section and click "Add Bot"
4. Copy the token and replace the current value in `discord.token`
5. Enable the necessary "Privileged Gateway Intents" (MESSAGE CONTENT, SERVER MEMBERS, PRESENCE)
6. In the "OAuth2" section, use the URL Generator to create an invite link with bot permissions and select the necessary permissions

### Administrators:

- Add the user IDs who will be bot administrators to the `discord.administrators` array
- Example: `"administrators": ["123456789012345678", "876543210987654321"]`

### Channels:

- Add the channel IDs where the bot can operate to the `discord.channel_ids` array
- To get IDs, enable "Developer Mode" in Discord (Settings > Advanced)
- Right-click on channels/users and select "Copy ID"

## 2. ZeroTier VPN Configuration

In the `vpn_handler_config.toml` file, you need to configure:

### ZeroTier API Key:

1. Create an account at [my.zerotier.com](https://my.zerotier.com)
2. Go to "Account" and generate an API Access Token
3. Add the token to the `api_key` field in the `[vpn_config]` section

### Network ID:

1. Create a new network in the ZeroTier panel
2. Copy the Network ID (format: 8056c2e21c000001)
3. Add it to the `network_id` field in the `[vpn_config]` section
4. Configure network access permissions as needed

## 3. Minecraft Server Configuration

### Minecraft Environment:

The system supports both **Fabric** and **Forge** mod loaders. You can configure the Minecraft server by editing the `.env` file in the respective server folder:

**For Fabric** (`minecraft_server/fabric/.env`):
```env
NETWORK_ID=your_zerotier_network_id
```

**For Forge** (`minecraft_server/forge/.env`):
```env
NETWORK_ID=your_zerotier_network_id
```

### Choosing the Server Version

The Minecraft version, mod loader version, Java version, and other settings are defined in the Docker image tag used in docker-compose.yml.
For example:
```
image: ghcr.io/gustamantovani/admine/minecraft_server:mc-1.21.7-forge57.0.2-java24-graalvm-zerotier
```

To change the server version, update the image tag accordingly. For instance, to use Fabric 0.17.2 on Minecraft 1.21.1:

```
image: ghcr.io/gustamantovani/admine/minecraft_server:mc-1.21.1-fabric0.17.2-installer1.1.0-java21-graalvm-zerotier
```

### Mods Configuration:

- Place your mod files (`.jar`) in the `mods/` folder of your chosen server type:
  - **Fabric mods**: `minecraft_server/fabric/mods/`
  - **Forge mods**: `minecraft_server/forge/mods/`
- Ensure mods are compatible with your chosen Minecraft version
- The server will automatically load mods from this folder on startup

### Server Properties:

- Server configurations are located in the `config/` folder
- You can modify `server.properties`, `eula.txt`, and other server settings

#### `server.properties` Configuration:

This file contains the main server settings. Key configurations include:

- **Server Details**:
  - `server-port=25565` - Port for players to connect
  - `rcon.port=25575` - RCON port for remote administration
  - `rcon.password=password` - RCON password (change this!)
  - `motd=Your Server Message` - Message displayed in server list

- **Gameplay Settings**:
  - `difficulty=normal` - Game difficulty (peaceful, easy, normal, hard)
  - `gamemode=survival` - Default game mode
  - `max-players=10` - Maximum number of players
  - `allow-flight=true` - Allow players to fly
  - `enable-command-block=true` - Enable command blocks

- **Security**:
  - `enforce-whitelist=true` - Only whitelisted players can join
  - `online-mode=false` - Set to true for premium servers
  - `rcon.password=password` - **⚠️ IMPORTANT: Change this password!**

### 🔒 RCON Password Security:

**It is highly recommended to change the default RCON password** for security reasons. You need to update it in **both** locations:

1. **Server Properties** (`minecraft_server/fabric/config/server.properties` or `minecraft_server/forge/config/server.properties`):
   ```properties
   rcon.password=your_secure_password
   ```

2. **Server Handler Configuration** (`server_handler/server_handler_config.yaml`):
   ```yaml
   minecraft_server:
     docker:
       rcon_password: "your_secure_password"
   ```

**Make sure both passwords match exactly**, otherwise the server handler won't be able to communicate with the Minecraft server.

#### `user_jvm_args.txt` Configuration:

This file controls Java Virtual Machine settings for server performance:

```plaintext
# Memory allocation (adjust based on your server's RAM)
-Xms2G    # Minimum RAM allocation (2GB)
-Xmx2G    # Maximum RAM allocation (2GB)

# For better performance, you can add:
# -XX:+UseG1GC                    # Use G1 garbage collector
# -XX:+UnlockExperimentalVMOptions
# -XX:MaxGCPauseMillis=100        # Reduce lag spikes
```

**Memory Recommendations**:
- **Light modpacks**: 2-4GB (`-Xmx4G`)
- **Medium modpacks**: 4-6GB (`-Xmx6G`)
- **Heavy modpacks**: 6-8GB+ (`-Xmx8G`)

## 4. System Initialization

After configuring your credentials, you can start the system with:

```bash
./start.sh
```

## System Usage

- Bot commands are available as slash commands (/) in Discord
- After starting the system, commands will be automatically registered in Discord
- To manage the server or VPN, use the available commands in Discord

## Functionality Verification

- Check if the bot is online in Discord
- Verify if slash commands are available on the server
- Check if Docker containers are running with `docker ps`

Now your Admine should be working correctly!

## About the Project

For more detailed information about the full Admine project, including documentation, architecture details, and development guides, please visit the main repository:

**🔗 [Admine - Full Project Repository](https://github.com/GustaMantovani/Admine)**

This slim version contains only the essential components needed for deployment. The main repository includes additional features, development tools, and comprehensive documentation.