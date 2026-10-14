import { GameServerType } from "./game-server-type";

export const DEFAULT_GAME_PORTS: Record<GameServerType, number> = {
	[GameServerType.Java]: 25565,
	[GameServerType.Bedrock]: 19132,
	[GameServerType.Hytale]: 5520,
};

export const DEFAULT_RCON_PORTS: Record<
	Exclude<GameServerType, GameServerType.Bedrock | GameServerType.Hytale>,
	number
> = {
	[GameServerType.Java]: 25575,
};

export const GAME_PORT_RANGES: Record<
	GameServerType,
	{ min: number; max: number }
> = {
	[GameServerType.Java]: { min: 25565, max: 25999 },
	[GameServerType.Bedrock]: { min: 19132, max: 19999 },
	[GameServerType.Hytale]: { min: 55200, max: 55999 },
};

export const RCON_PORT_RANGES: Record<
	Exclude<GameServerType, GameServerType.Bedrock | GameServerType.Hytale>,
	{ min: number; max: number }
> = {
	[GameServerType.Java]: { min: 26000, max: 26999 },
};

export function resolveServerPort(serverType: GameServerType): number {
	return DEFAULT_GAME_PORTS[serverType];
}
