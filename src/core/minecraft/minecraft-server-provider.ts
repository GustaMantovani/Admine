import type { CreateMinecraftServerRequest } from "../types/create-minecraft-server-request";
import type { MinecraftServer } from "../types/minecraft-server";
import type { MinecraftServerStats } from "../types/minecraft-server-stats";
import type { ProviderName } from "../types/provider-name";

export interface MinecraftServerProvider {
	readonly name: ProviderName;
	create(
		createServerRequest: CreateMinecraftServerRequest,
		acceptEula?: boolean,
	): Promise<string>;
	list(): Promise<MinecraftServer[]>;
	info(id: string): Promise<MinecraftServer>;
	stats(id: string): Promise<MinecraftServerStats>;
	logs(id: string): Promise<string[]>;
	start(id: string): Promise<void>;
	stop(id: string): Promise<void>;
	restart(id: string): Promise<void>;
	kill(id: string): Promise<void>;
	backup(id: string): Promise<void>;
	update(id: string): Promise<void>;
	clone(id: string): Promise<void>;
	acceptEula(id: string): Promise<void>;
	sendCommand(id: string, command: string): Promise<void>;
	delete(id: string, removeFiles: boolean): Promise<void>;
	stopAll(): Promise<number>;
}
