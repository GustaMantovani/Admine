import type { GameServerType } from "./game-server-type";
import type { MinecraftServerUserDefinedProperties } from "./minecraft-server-user-defined-properties";

export type CreateGameServerRequest = {
	name: string;
	type: GameServerType;
	version: string;
	minimumMemory: number;
	maximumMemory: number;
	port: number;
	rconPort?: number;
	properties?: MinecraftServerUserDefinedProperties;
};
