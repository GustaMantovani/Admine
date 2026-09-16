import { MinecraftServerType } from "./domain/minecraft-server-type";

export const SERVERS_HOST_FOR_ADMINE = process.env.ADMINE_SERVERS_HOST_FOR_ADMINE || '127.0.0.1';

// Defaults
export const DEFAULT_MINECRAFT_SERVER_TYPE = MinecraftServerType.Java;
export const DEFAULT_MINECRAFT_SERVER_VERSION = '26.3';
export const DEFAULT_JAVA_PORT = 25565;
export const DEFAULT_BEDROCK_PORT = 19132;
export const DEFAULT_HYTALE_PORT = 5520;
export const DEFAULT_MINIMUM_MEMORY = 1;
export const DEFAULT_MAXIMUM_MEMORY = 2;

export const DEFAULT_CRAFTY_API_BASE_PATH = 'https://localhost:8443';
