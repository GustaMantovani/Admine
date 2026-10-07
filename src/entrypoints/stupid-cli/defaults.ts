import { MinecraftServerType } from "../../core/types/minecraft-server-type";
import type { MinecraftServerWebhookUrls } from "../../core/types/minecraft-server-webhook-trigger";

// Files
export const DEFAULT_CLI_CONFIG_FILE_PATH = "./config/cli-config.toml";
export const DEFAULT_CRAFTY_CREDS_FILE_PATH =
	".data/crafty/config/default-creds.txt";

// Minecraft servers
export const DEFAULT_SERVERS_HOST_FOR_ADMINE = "127.0.0.1";
export const DEFAULT_MINECRAFT_SERVER_TYPE = MinecraftServerType.Java;
export const DEFAULT_MINECRAFT_SERVER_VERSION = "26.3";
export const DEFAULT_MINECRAFT_SERVER_MINIMUM_MEMORY = 2;
export const DEFAULT_MINECRAFT_SERVER_MAXIMUM_MEMORY = 4;
export const DEFAULT_MINECRAFT_SERVER_RCON_PORT = "25575";
export const DEFAULT_MINECRAFT_SERVER_RCON_PASSWORD = "admine";
export const DEFAULT_MINECRAFT_SERVER_WHITE_LIST_ENABLE = false;

// Crafty
export const DEFAULT_CRAFTY_WEBHOOK_URLS: MinecraftServerWebhookUrls = {};
export const DEFAULT_CRAFTY_API_BASE_URL = "https://localhost:8443/api/v2";
export const DEFAULT_CRAFTY_INFRA_SETUP_COMMAND =
	'docker start crafty_container 2>/dev/null || docker run --name crafty_container --detach -p 8443:8443 -p 8123:8123 -p 19132:19132/udp -p 25500-25600:25500-25600 -e TZ=Etc/UTC -v "/$(pwd)/.data/crafty/backups:/crafty/backups" -v "/$(pwd)/.data/crafty/logs:/crafty/logs" -v "/$(pwd)/.data/crafty/servers:/crafty/servers" -v "/$(pwd)/.data/crafty/config:/crafty/app/config" -v "/$(pwd)/.data/crafty/import:/crafty/import" arcadiatechnology/crafty-4:latest';
// This command must return a success exit code when the server is running, or any other exit code otherwise
export const DEFAULT_CRAFTY_INFRA_CHECK_COMMAND =
	"docker inspect -f '{{.State.Running}}' crafty_container | grep -q true";
export const DEFAULT_CRAFTY_INFRA_CHECK_DELAY_MS = 10000;
export const DEFAULT_CRAFTY_WEBHOOKS_ENABLED = false;