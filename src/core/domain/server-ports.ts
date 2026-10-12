import type { CreateGameServerRequest } from "./create-game-server-request";
import { GameServerType } from "./game-server-type";

// Default ports are defined by each game's protocol
const DEFAULT_PORTS: Record<GameServerType, number> = {
	[GameServerType.Java]: 25565,
	[GameServerType.Bedrock]: 19132,
	[GameServerType.Hytale]: 5520,
};

export function resolveServerPort(request: CreateGameServerRequest): number {
	return request.port ?? DEFAULT_PORTS[request.type] ?? DEFAULT_PORTS.java;
}
