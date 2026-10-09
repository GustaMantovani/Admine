import type { MinecraftServerType } from "../types/minecraft-server-type";
import type { MinecraftServerUserDefinedProperties } from "./minecraft-server-user-defined-properties";

export type CreateMinecraftServerRequest = {
	name: string;
	type: MinecraftServerType;
	version: string;
	minimumMemory: number;
	maximumMemory: number;
	// Falls back to the protocol's default port for the server type
	port?: number;
	properties?: MinecraftServerUserDefinedProperties;
};
