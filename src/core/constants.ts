import { MinecraftServerType } from "./types/minecraft-server-type";

export const SERVERS_HOST_FOR_ADMINE =
	Bun.env.ADMINE_SERVERS_HOST_FOR_ADMINE || "127.0.0.1";

// Defaults
export const DEFAULT_MINECRAFT_SERVER_TYPE = (Bun.env
	.ADMINE_DEFAULT_MINECRAFT_SERVER_TYPE ||
	MinecraftServerType.Java) as MinecraftServerType;
export const DEFAULT_MINECRAFT_SERVER_VERSION =
	Bun.env.ADMINE_DEFAULT_MINECRAFT_SERVER_VERSION || "26.3";
export const DEFAULT_JAVA_PORT = Number(
	Bun.env.ADMINE_DEFAULT_JAVA_PORT || 25565,
);
export const DEFAULT_BEDROCK_PORT = Number(
	Bun.env.ADMINE_DEFAULT_BEDROCK_PORT || 19132,
);
export const DEFAULT_HYTALE_PORT = Number(
	Bun.env.ADMINE_DEFAULT_HYTALE_PORT || 5520,
);
export const DEFAULT_MINIMUM_MEMORY = Number(
	Bun.env.ADMINE_DEFAULT_MINIMUM_MEMORY || 2,
);
export const DEFAULT_MAXIMUM_MEMORY = Number(
	Bun.env.ADMINE_DEFAULT_MAXIMUM_MEMORY || 4,
);

// Crafty
export const CRAFTY_API_BASE_PATH =
	Bun.env.ADMINE_CRAFTY_API_BASE_PATH || "https://localhost:8443/api/v2";
//export const CRAFTY_API_TOKEN = requiredEnv("ADMINE_CRAFTY_API_TOKEN");
export const CRAFTY_API_TOKEN = Bun.env.ADMINE_CRAFTY_API_TOKEN;
export const CRAFTY_SELF_HOSTED = Bun.env.ADMINE_CRAFTY_SELF_HOSTED === "true";
export const CRAFTY_INFRAESTRUTURE_COMMAND_STRING_SETUP =
	Bun.env.ADMINE_CRAFTY_INFRAESTRUTURE_COMMAND_STRING_SETUP ||
	'docker run --name crafty_container --detach -p 8443:8443 -p 8123:8123 -p 19132:19132/udp -p 25500-25600:25500-25600 -e TZ=Etc/UTC -v "/$(pwd)/crafty-data/backups:/crafty/backups" -v "/$(pwd)/crafty-data/logs:/crafty/logs" -v "/$(pwd)/crafty-data/servers:/crafty/servers" -v "/$(pwd)/crafty-data/config:/crafty/app/config" -v "/$(pwd)/crafty-data/import:/crafty/import" arcadiatechnology/crafty-4:latest';
export const CRAFTY_INFRAESTRUTURE_CHECK_DELAY_MS = Number(
	Bun.env.ADMINE_CRAFTY_INFRAESTRUTURE_CHECK_DELAY_MS || 10000,
);
// This command must return a success exit code when the server is running, or any other exit code otherwise
export const CRAFTY_INFRAESTRUTURE_COMMAND_SRTING_CHECK =
	Bun.env.ADMINE_CRAFTY_INFRAESTRUTURE_COMMAND_SRTING_CHECK ||
	"docker inspect -f '{{.State.Running}}' crafty_container | grep -q true";
