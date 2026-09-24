import type { MinecraftServerType } from "../types/minecraft-server-type";

export type CreateMinecraftServerRequest = {
	name: string;
	type?: MinecraftServerType;
	version?: string;
	minimumMemory?: number;
	maximumMemory?: number;
	port?: number;
};
