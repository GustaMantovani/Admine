import type { CreateMinecraftServerRequest } from "../types/create-minecraft-server-request";
import { MinecraftServerType } from "../types/minecraft-server-type";

// Default ports are defined by each game's protocol
const DEFAULT_PORTS: Record<MinecraftServerType, number> = {
	[MinecraftServerType.Java]: 25565,
	[MinecraftServerType.Bedrock]: 19132,
	[MinecraftServerType.Hytale]: 5520,
};

export function resolveServerPort(
	request: CreateMinecraftServerRequest,
): number {
	return request.port ?? DEFAULT_PORTS[request.type] ?? DEFAULT_PORTS.java;
}
